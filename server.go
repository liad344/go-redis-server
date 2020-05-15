package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/evio"
	"go.uber.org/zap"
	"net"
	"strings"
	"sync"
)

type Mux map[string]Handler
type FastMux map[string]FastHandler
type key string

type Conn struct {
	evio.Conn
}
type Server struct {
	cfg    ServerCfg
	conns  map[*Conn]bool
	ins    *RedisInstance
	mux    *Mux
	fmux   *FastMux
	ln     net.Listener
	logger *zap.SugaredLogger
}

func (s *Server) Init() {
	s.fmux = NewFastMux()
	s.ins = NewInstance()
	s.initConfig()
	s.fmux.FHandleFunc("set", s.ins.FSet)
	s.fmux.FHandleFunc("get", s.ins.FGet)
}

func NewFastMux() *FastMux {
	m := make(FastMux)
	return &m
}

func NewServer() *Server {
	return &Server{
		cfg:   ServerCfg{},
		conns: make(map[*Conn]bool),
		ln:    nil,
		mux:   &Mux{},
	}
}
func (m FastMux) FHandleFunc(cmd string, h FastHandler) {
	if _, ok := m[cmd]; ok {
		log.Error("Handler already exist")
	}
	m[cmd] = h
}

func (m Mux) HandleFunc(cmd string, h Handler) {
	if _, ok := m[cmd]; ok {
		log.Error("Handler already exist")
	}
	m[cmd] = h
}

func NewInstance() *RedisInstance {
	return &RedisInstance{
		data:  map[key]value{},
		Mutex: sync.Mutex{},
	}
}

func NewMux() *Mux {
	m := make(Mux)
	return &m
}

var i int

func (s *Server) FastListenAndServerRESP() {
	var events evio.Events
	m := *s.fmux
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		cmd := parseCmd(in)
		con := Conn{c}
		if h, ok := m[strings.ToLower(string(cmd.Args[0]))]; ok {
			out = h(con, cmd)
		} else {
			log.Error("No handler for ", string(cmd.Args[0]), " command")
		}
		return out, evio.Close
	}
	log.Info("Serving Fast?")
	events.NumLoops = 3
	if err := evio.Serve(events, s.cfg.Addr); err != nil {
		panic(err.Error())
	}

}
