package main

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"sync"
	"time"
)

type serverCfg struct {
	addr string
}

type Handler func (conn Conn, cmd Command)
type key string

type handleConnection func(conn net.Conn) bool
type handleClosedConnection func(conn Conn, err error)
type Mux map[string]Handler

type Conn struct {
	net.Conn
}
type Server struct {
	cfg serverCfg
	conns   map[*net.Conn]bool
	ins    *Instance
	mux    *Mux
	accept handleConnection
	closed handleClosedConnection

	ln      net.Listener
}

func (s *Server) Init(){
	s.mux = NewMux()
	s.ins = NewInstance()
	s.accept = onNewConnection
	s.closed = onConnectionClosed
	s.mux.HandleFunc("set" , s.ins.Set)
	s.mux.HandleFunc("get" , s.ins.Get)
	s.mux.HandleFunc("del" , s.ins.Del)
	s.mux.HandleFunc("ping" , s.ins.Ping)
}
func (m Mux) HandleFunc(cmd string, h Handler) {
	if  _ , ok := m[cmd]; ok {
		log.Error("Handler already exist")
	}
	m[cmd] = h
}

func NewInstance() *Instance {
	
}

func NewMux() *Mux {
	m := make(Mux)
	return &m
}
func (s *Server) ListenAndServerRESP(){
	log.Info("Serving connections")
	for {
		conn , err := s.ln.Accept()
		if err != nil {
			log.Error("Could not connect")
		}
		s.conns[&conn] = true
		go handle(conn , *s.mux)


	}
}

func handle(conn net.Conn , mux Mux) {
	buf , err := readCmd(conn)
	if err != nil {
		log.Error("Could not read form connection " , err)
		return
	}
	log.Info(buf)
	cmd := parseCmd(buf)
	for _ ,a := range cmd.Args{
		log.Info(string(a))
	}
	if h , ok := mux[string(cmd.Args[0])]; ok {
		h(conn , cmd)
	}else {
		log.Error("No handler for " , string(cmd.Args[0]) , " command")
		return
	}

}

func readCmd(conn net.Conn) ([]byte , error) {
	return ioutil.ReadAll(conn)
}