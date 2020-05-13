package main

import (
	log "github.com/sirupsen/logrus"
	"net"
	"regexp"
	"strconv"
	_ "strconv"
)

 func parseCmd(buf []byte) (C Command) {
	C.Data = buf
	 r, ArgsRegex := regex()
	 if loc := r.FindIndex(buf); loc != nil {
		 makeArgs(buf, loc, C)
		 if argsIndexes := ArgsRegex.FindAllIndex(buf[loc[1]:], -1); argsIndexes != nil {
			 for _, index := range argsIndexes {
				 start := loc[1] + index[0]
				 end := loc[1] + index[1]
				 //log.Info("Searching ", string(buf[loc[1]:]))
				 //log.Info("index ", index, " val ", string(buf[start:end]))
				 C.Args = append(C.Args, buf[start+2:end-2])
			 }
		 }
	 }
	 return C
 }

func makeArgs(buf []byte, loc []int, C Command) {
	reqLen, err := strconv.ParseInt(string(buf[loc[0]+1]), 10, 8)
	if err != nil {
		log.Error("Could not parse cliet request len ", err)
		return
	}
	C.Args = make([][]byte, reqLen)
}

func regex() (*regexp.Regexp, *regexp.Regexp) {
	r, err := regexp.Compile("\\*[1-9]\\r\\n")
	if err != nil {
		log.Error("Could not compile regex ", err)
		return nil, nil
	}
	ArgsRegex, err := regexp.Compile("\\r\\n.*\\r\\n")
	if err != nil {
		log.Error("Could not compile regex ", err)
		return nil, nil
	}
	return r, ArgsRegex
}

func readCmd(conn net.Conn) (b []byte  ,err error) {
	b =  make([]byte , 64)
	_ , err = conn.Read(b)
	//log.Info(string(b) , " from ip " , conn.RemoteAddr())
	return b , err
}

