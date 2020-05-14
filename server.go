package main

import (
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"io"
	"net"
	"runtime"
	"strings"
	"sync"
)

type Mux map[string]Handler
type key string

type Conn struct {
	net.Conn
}
type Server struct {
	cfg   ServerCfg
	conns map[*net.Conn]bool
	ins   *RedisInstance
	mux   *Mux

	ln     net.Listener
	logger *zap.SugaredLogger
}

func (s *Server) Init() {
	s.mux = NewMux()
	s.ins = NewInstance()
	s.initConfig()
	log.Info("Using config: ", s.cfg)
	s.mux.HandleFunc("set", s.ins.Set)
	s.mux.HandleFunc("get", s.ins.Get)
	s.mux.HandleFunc("del", s.ins.Del)
	s.mux.HandleFunc("ping", s.ins.Ping)
	ln, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		log.Error(err)
	}
	s.ln = ln
	//logger, _ := zap.NewProduction()
	//s.logger = logger.Sugar()

}

func NewServer() *Server {
	return &Server{
		cfg:   ServerCfg{},
		conns: make(map[*net.Conn]bool),
		ln:    nil,
		mux:   &Mux{},
	}
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

func (s *Server) ListenAndServerRESP() {
	log.Info("Serving connections")
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Error("Could not connect")
		}
		if i++; i <= runtime.GOMAXPROCS(s.cfg.MaxGoRoutines) {
			go handleClient(conn, *s.mux)
			i--
		}
	}
}

func handleClient(conn net.Conn, mux Mux) {
	for {
		if !handle(conn, mux) {
			log.Info("Stopping connection w/ ", conn.RemoteAddr())
			break
		}
	}
	err := conn.Close()
	if err != nil {
		log.Error("Could not close connection ", err)
	}
}
func handle(conn net.Conn, mux Mux) bool {
	buff, err := readCmd(conn)
	if err != nil && err != io.EOF {
		log.Error("Could  not read form connection ", err)
		return false
	}
	if len(buff) == 0 || err == io.EOF {
		return false
	}

	cmd := parseCmd(buff)
	c := Conn{conn}
	if h, ok := mux[strings.ToLower(string(cmd.Args[0]))]; ok {
		h(c, cmd)
	} else {
		log.Error("No handler for ", string(cmd.Args[0]), " command")
		return false
	}
	return true
}
