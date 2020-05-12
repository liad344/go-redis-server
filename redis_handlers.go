package main

import (
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)



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


func (i *Instance) Ping(conn Conn, cmd Command) {
	conn.WriteString("PONG")
	log.Info("Ponged ip" , conn.RemoteAddr() )
}


func (i *Instance) Del(conn Conn, cmd Command) {
	if len(cmd.Args) != 2 {
		conn.WriteString("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	i.Lock()
	//del
	i.Unlock()
	if !ok {
		conn.WriteString(0)
	} else {
		conn.WriteString(1)
	}
	log.Info("Deleted")
	return
}

func (i *Instance) Get(conn Conn, cmd Command) {
	if len(cmd.Args) < 2 {
		conn.WriteString("Not enough arguments")
		return
	}
	i.Lock()
	//get
	i.Unlock()
	if !ok {
		conn.WriteNull()
	} else {
		conn.WriteBulk(val.data)
	}
	log.Info("Got val " , string(val))
	return
}

func (i *Instance) Set(conn Conn, cmd Command) {
	if len(cmd.Args) < 3 {
		conn.WriteString("Not enough arguments")
		return
	}
	i.Lock()
	//Set
	i.Unlock()
	conn.WriteString("OK")
	log.Info("Set " , string(cmd.Args[2]))
	return
}

