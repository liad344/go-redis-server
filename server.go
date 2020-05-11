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

type Handler func (conn net.Conn, cmd Command)
type key string

type handleConnection func(conn net.Conn) bool
type handleClosedConnection func(conn net.Conn, err error)
type Mux map[string]Handler

type Command struct {
	Data []byte
	Args [][]byte
}

type value struct {
	data []byte
	ttl time.Duration
}

type Instance struct {
	data map[key]value
	sync.Mutex
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