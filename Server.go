package main

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"regexp"
	"strconv"
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
	cmd := parseCmd(buf)
	if h , ok := mux[string(cmd.Args[0])]; ok {
		h(conn , cmd)
	}else {
		log.Error("No handler for " , string(cmd.Args[0]) , " command")
		return
	}

}
	// *2\r\n$3\r\nget\r\n$2\r\nyo\r\n\ get (yo)
//*1\r\n$4\r\nping\r\n\
//*5\r\n$3\r\nset\r\n$2\r\nyo\r\n$2\r\nyo\r\n$2\r\nex\r\n$1\r\n1\r\n\ set(yo , yo , 1)
func parseCmd(buf []byte)(C Command) {
	C.Data = buf
	//todo error handling here is gay AF
	r , _ := regexp.Compile("\\*[1-9]\\r\\n\\$")
	if loc := r.FindIndex(buf); loc != nil {
		log.Info(loc[0]+1)
		reqLen , _  := strconv.ParseInt(string(buf[loc[0]+1]) , 10 , 8)
		C.Args = make([][]byte , reqLen)
		ArgsRegex , _ :=  regexp.Compile("\\r\\n.*\\r\\n")
		if argsIndexes := ArgsRegex.FindAllIndex(buf[loc[1]:] , -1); argsIndexes != nil {
			for _ , index := range argsIndexes {
				//log.Info("loc[1] " , loc[1] , " loc[0] ", loc[0])
				start := loc[1] + index[0]
				end := loc[1] + index[1]
			//log.Info("Searching " , string(buf[loc[1]:]))
			//	log.Info("index " , index , " val " , string(buf[start:end]))
				C.Args = append(C.Args, buf[start+2:end-2])
			}
		}
	}
	return C
}

func readCmd(conn net.Conn) ([]byte , error) {
	return ioutil.ReadAll(conn)
}