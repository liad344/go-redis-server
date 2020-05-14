package main

import (
	"bytes"
	valid "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"net"
	"regexp"
	"strconv"
	_ "strconv"
)

const BuffSize = 64

func parseCmd(buff []byte) (C Command) {
	C.Data = buff
	//_ , _ := regex()
	//log.Info("buff " , buff)
	//log.Info("buff str " , string(buff))
	if loc := fastIndex(buff); loc != nil {
		makeArgs(buff, loc, C)
		if argsIndexes := fastArgsIndex(buff[loc[1]:], 3); argsIndexes != nil {
			for _, index := range argsIndexes {
				start := loc[1] + index[0]
				end := loc[1] + index[1]
				//log.Info("Searching ", string(buff[loc[1]:]))
				//log.Info("index ", index, " val ", string(buff[start:end]))
				C.Args = append(C.Args, buff[start:end])
			}
		}
	}
	return C
}

func fastIndex(buff []byte) []int {
	if buff[0] != '*' {
		return nil
	}
	if !valid.IsInt(string(buff[1])) {
		return nil
	}
	if buff[2] != '\r' && buff[3] != 'n' && buff[4] != '$' {
		return nil
	}

	return []int{0, 4}
}

func fastArgsIndex(buff []byte, l int) [][]int {
	indexes := make([][]int, l)
	crlf := []byte{'\r', '\n'}
	j := 0
	i := 0
	var n2 int
	for {
		firstCrlf := bytes.Index(buff[i:], crlf)
		if firstCrlf == -1 {
			break
		}
		//log.Info("buff str at fastArgIndex " , string(buff[i:]))
		//log.Info("firstCrlf " , firstCrlf)
		if n, err := strconv.Atoi(string(buff[i+1 : i+firstCrlf])); err == nil && buff[i] == '$' {
			//log.Info("n " , n)
			//log.Info("strconv " , string(buff[i+1: i + firstCrlf]))
			indexes[j] = []int{i + firstCrlf + 2, i + firstCrlf + 2 + n}
			//log.Info("Args #" , j , " " , indexes[j] , " Value " , string(buff[i+ firstCrlf+2:i+firstCrlf+2+n]))
			n2 = n
		} else {
			log.Error(err)
		}
		i += n2 + firstCrlf + 2
		lastCrlf := bytes.Index(buff[i:], crlf)
		i += lastCrlf + 2
		//log.Info("i " , i , " buff[i] " , string(buff[i]))
		j++
		//log.Info("len buff " , len(buff) )
		//log.Info("j " , j , " l " , l)
		if i >= len(buff) || buff[i] != '$' || j == l {
			break
		}
	}
	//log.Info("RETURNING " , indexes)
	return indexes
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

func readCmd(conn net.Conn) (b []byte, err error) {
	buff := bytes.NewBuffer(make([]byte, 0, 1024))
	for {
		b = make([]byte, BuffSize)
		n, err := conn.Read(b)
		if n < BuffSize && n != 0 {
			buff.Write(b)
			return buff.Bytes(), nil
		}
		if err != nil {
			return nil, err
		}
		buff.Write(b)
	}
}
