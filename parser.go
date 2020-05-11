package main

import (
	log "github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	_ "strconv"
)

const (
	INT = ":"
	SIMPLE_STRING = "+"
	ERRORS = "-"
	BULK_STRING = "$"
	ARRAY = "*"
)

// *2\r\n$3\r\nget\r\n$2\r\nyo\r\n\ get (yo)
//*1\r\n$4\r\nping\r\n\
//*5\r\n$3\r\nset\r\n$2\r\nyo\r\n$2\r\nyo\r\n$2\r\nex\r\n$1\r\n1\r\n\ set(yo , yo , 1)
func parseCmd(buf []byte) (C Command) {
	C.Data = buf
	r, err := regexp.Compile("\\*[1-9]\\r\\n\\$")
	if err != nil {
		log.Error("Could not compile regex " , err)
	}
	if loc := r.FindIndex(buf); loc != nil {
		reqLen, err := strconv.ParseInt(string(buf[loc[0]+1]), 10, 8)
		if err != nil {
			log.Error("Could not parse cliet request len " , err)
		}
		C.Args = make([][]byte, reqLen)
		ArgsRegex, err := regexp.Compile("\\r\\n.*\\r\\n")
		if err != nil {
			log.Error("Could not compile regex " , err)
		}
		if argsIndexes := ArgsRegex.FindAllIndex(buf[loc[1]:], -1); argsIndexes != nil {
			for _, index := range argsIndexes {
				start := loc[1] + index[0]
				end := loc[1] + index[1]
				C.Args = append(C.Args, buf[start+2:end-2])
			}
		}
	}
	return C
}

