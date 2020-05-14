package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"net"
	"regexp"
	"strconv"
	_ "strconv"
)

const BuffSize = 64

func parseCmd(buff []byte) (C Command) {
	C.Data = buff
	if loc := fastIndex(buff); loc != nil {
		numArgs, err := strconv.ParseInt(string(buff[loc[0]+1]), 10, 8)
		if err != nil {
			log.Error("Could not parse client request len ", err)
			return
		}
		C.Args = make([][]byte, numArgs)
		if argsIndexes := fastArgsIndex(buff[loc[1]:], int(numArgs)); argsIndexes != nil {
			for _, index := range argsIndexes {
				start := loc[1] + index[0]
				end := loc[1] + index[1]
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

func fastArgsIndex(buff []byte, numArgs int) [][]int {
	indexes := make([][]int, numArgs)
	crlf := []byte{'\r', '\n'}
	crlfLen := len(crlf)
	j := 0
	i := 0
	for {
		firstCrlf := bytes.Index(buff[i:], crlf)
		if firstCrlf == -1 {
			break
		}
		if n, err := strconv.Atoi(string(buff[i+1 : i+firstCrlf])); err == nil && buff[i] == '$' {
			indexes[j] = []int{i + firstCrlf + 2, i + firstCrlf + 2 + n}
			i += n + firstCrlf + crlfLen
		} else {
			log.Error(err)
			return nil
		}
		lastCrlf := bytes.Index(buff[i:], crlf)
		i += lastCrlf + crlfLen
		j++
		if /*i >= len(buff) || buff[i] != '$' || */ j == numArgs {
			break
		}
	}
	return indexes
}

func makeArgs(buf []byte, loc []int, C Command) {

}

func readCmd(conn net.Conn) (b []byte, err error) {
	buff := bytes.NewBuffer(make([]byte, 0, 256))
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
