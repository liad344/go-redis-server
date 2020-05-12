package main

import (
	"encoding/binary"
	log "github.com/sirupsen/logrus"
	//"net"
	//"strconv"
	"sync"
	"time"
)



type Command struct {
	Data []byte
	Args [][]byte
}

type value struct {
	data []byte
	//ttl time.Duration
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
	if len(cmd.Args) < 2 {
		conn.WriteString("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	i.Lock()
	//del
	i.Unlock()
	log.Info("Deleted")
	return
}

func (i *Instance) Get(conn Conn, cmd Command) {
	if len(cmd.Args) < 2 {
		conn.WriteString("Not enough arguments")
		return
	}
	// get yo
	key := key(cmd.Args[1])
	i.Lock()
	val , ok := i.data[key]
	i.Unlock()
	if !ok {
		conn.WriteNull()
	} else {
		conn.WriteString(string(val.data))
	}
	log.Info("Got val " , string(val.data))
	return
}

func (i *Instance) Set(conn Conn, cmd Command) {
	if len(cmd.Args) < 3 {
		conn.WriteString("Not enough arguments")
		return
	}
	val := value{data:cmd.Args[2]}
	key := key(cmd.Args[1])
	i.Lock()
	// args[0] = set , args[1] = key , args[2] = val , args[3] = ex , args[4] = time
	i.data[key] = val
	if string(cmd.Args[3]) == "ex" {
		i.ttl(key , cmd.Args[4])
	}
	i.Unlock()
	conn.WriteString("OK")
	log.Info("Set " , string(cmd.Args[2]))
	return
}

func (i *Instance) ttl(k key, t []byte) {
	tInt := binary.BigEndian.Uint64(t)
	<-time.After(time.Duration(tInt*1000000000))
	delete(i.data , k )
}

