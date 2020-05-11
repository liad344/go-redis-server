package main

import (
	log "github.com/sirupsen/logrus"
	"net"
)

func main(){
	s := &Server{
		cfg:   serverCfg{addr: ":8000"},
		conns: make(map[*net.Conn]bool),
		ln:    nil,
		mux: &Mux{},
	}

	ln , err := net.Listen("tcp" , ":8000")
	if err != nil {
		log.Error(err)
	}
	s.ln = ln
	s.ListenAndServerRESP()
}





