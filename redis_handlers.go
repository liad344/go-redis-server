package main

import (
	log "github.com/sirupsen/logrus"
	"strconv"

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
	//deleteAfterTtl time.Duration
}


type RedisInstance struct {
	data map[key]value
	sync.Mutex
}


func (i *RedisInstance) Ping(conn Conn, cmd Command) {
	conn.WriteString("PONG")
	log.Info("Ponged ip" , conn.RemoteAddr() )
}


func (i *RedisInstance) Del(conn Conn, cmd Command) {
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

func (i *RedisInstance) Get(conn Conn, cmd Command) {
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
		conn.WriteString("Key not found")
	} else {
		conn.WriteString(string(val.data))
	}
	log.Info("Got val " , string(val.data))
	return
}

func (i *RedisInstance) Set(conn Conn, cmd Command) {
	if len(cmd.Args) < 3 {
		conn.WriteString("Not enough arguments")
		return
	}
	val := value{data:cmd.Args[2]}
	key := key(cmd.Args[1])
	i.Lock()
	// args[0] = set , args[1] = key , args[2] = val , args[3] = ex , args[4] = time
	i.data[key] = val
	if string(cmd.Args[3]) == "ex" || string(cmd.Args[3]) == "px"{
		go i.deleteAfterTtl(key , cmd.Args[4])
	}
	i.Unlock()
	conn.WriteString("OK")
	log.Info("Set " , string(cmd.Args[2]))
	return
}

func (i *RedisInstance) deleteAfterTtl(k key, t []byte) {
	tInt , err := strconv.ParseInt(string(t) , 10 , 8)
	if err != nil{
		log.Error("Could not parse ttl")
	}
	log.Info("deleting after " , time.Second * time.Duration(tInt))
	<-time.After(time.Second * time.Duration(tInt))
	delete(i.data , k )
}

