package main

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

type Command struct {
	Data []byte
	Args [][]byte
}

type value struct {
	data      []byte
	redisType string
	//ttl time.Duration
}

type RedisInstance struct {
	data map[key]value
	sync.Mutex
}

func (i *RedisInstance) Ping(conn Conn, cmd Command) {
	conn.Write(STRING, []byte("PONG"))
	log.Info("Ponged ", conn.RemoteAddr())
}

func (i *RedisInstance) Del(conn Conn, cmd Command) {
	if len(cmd.Args) < 2 {
		conn.Write(ERROR, []byte("Wrong number of arguments"))
		return
	}
	i.Lock()
	delete(i.data, key(cmd.Args[1]))
	i.Unlock()
	log.Info("Deleted")
}

func (i *RedisInstance) Get(conn Conn, cmd Command) {
	if len(cmd.Args) < 2 {
		conn.Write(ERROR, []byte("Not enough arguments"))
		return
	}
	key := key(cmd.Args[1])
	i.Lock()
	val, ok := i.data[key]
	i.Unlock()
	if ok {
		conn.Write(val.redisType, val.data)
		conn.Close()
		return
	}

	conn.Write(ERROR, []byte("Key not found"))
	conn.Close()
}

func (i *RedisInstance) Set(conn Conn, cmd Command) {
	if len(cmd.Args) < 3 {
		conn.Write(ERROR, []byte("Not enough arguments"))
		return
	}
	val := value{data: cmd.Args[2], redisType: STRING}
	key := key(cmd.Args[1])

	safeSet(i, key, val, cmd)

	conn.Write(STRING, []byte("OK"))
	conn.Close()
	return
}

func safeSet(i *RedisInstance, key key, val value, cmd Command) {
	i.Lock()
	i.data[key] = val
	if len(cmd.Args) > 3 {
		go i.deleteAfterTtl(key, cmd.Args[3], cmd.Args[4])
	}
	i.Unlock()
}

func (i *RedisInstance) deleteAfterTtl(k key, ttl []byte, timeFormat []byte) {
	if string(timeFormat) != "px" || string(timeFormat) != "ex" {
		return
	}
	tInt, err := strconv.ParseInt(string(ttl), 10, 8)
	if err != nil {
		log.Error("Could not parse ttl")
	}
	<-time.After(time.Second * time.Duration(tInt))
	delete(i.data, k)
}
