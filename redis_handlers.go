package main

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

type Handler func(conn Conn, cmd Command)
type FastHandler func(conn Conn, cmd Command) []byte

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

func (i *RedisInstance) FGet(conn Conn, cmd Command) []byte {
	if len(cmd.Args) < 2 {
		return conn.FWrite(ERROR, []byte("Not enough arguments"))
	}
	key := key(cmd.Args[1])
	i.Lock()
	val, ok := i.data[key]
	i.Unlock()
	if ok {
		return conn.FWrite(val.redisType, val.data)
	}

	return []byte("Key not found")
}

func (i *RedisInstance) FSet(conn Conn, cmd Command) []byte {
	if len(cmd.Args) < 3 {
		return conn.FWrite(ERROR, []byte("Not enough arguments"))
	}
	val := value{data: cmd.Args[2], redisType: STRING}
	key := key(cmd.Args[1])

	safeSet(i, key, val, cmd)

	return conn.FWrite(STRING, []byte("OK"))
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
