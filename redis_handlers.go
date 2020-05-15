package main

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

type Handler func(conn Conn, cmd Command, i *RedisInstance)

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
type RedisMaster struct {
	ins []*RedisInstance
	RedisInstance
}

func (m *RedisMaster) NewInstance() *RedisInstance {
	newIns := &RedisInstance{
		data:  m.data,
		Mutex: sync.Mutex{},
	}
	m.ins = append(m.ins, newIns)
	return newIns
}

func (m *RedisMaster) Get(conn Conn, cmd Command, i *RedisInstance) {
	if len(cmd.Args) < 2 {
		conn.Write(ERROR, []byte("Not enough arguments"))
		return
	}
	key := key(cmd.Args[1])
	if get(conn, i, key) {
		log.Info("Got before sync")
		return
	}
	log.Info("SyncFromMaster")
	m.syncFromMaster(i)

	if get(conn, i, key) {
		log.Info("Got After sync")
		return
	}

	conn.Write(ERROR, []byte("Key not found"))
}

func get(conn Conn, i *RedisInstance, key key) bool {
	val, ok := i.data[key]
	if ok {
		conn.Write(val.redisType, val.data)
		return true
	}
	return false
}

func (m *RedisMaster) Set(conn Conn, cmd Command, i *RedisInstance) {
	l := len(cmd.Args)
	if l < 3 {
		conn.Write(ERROR, []byte("Not enough arguments"))
		return
	}
	val := value{data: cmd.Args[2], redisType: STRING}
	key := key(cmd.Args[1])

	i.data[key] = val
	conn.Write(STRING, []byte("OK"))
	if l > 3 {
		go i.deleteAfterTtl(key, cmd.Args[3], cmd.Args[4])
	}
	return
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
