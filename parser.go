package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
	_ "strconv"
)

const BuffSize = 64

func parseCmd(buff []byte) (C Command) {
	C.Data = buff
	if validRequest(buff) {
		numArgs, err := strconv.ParseInt(string(buff[1]), 10, 8)
		if err != nil {
			log.Error("Could not parse client request len ", err)
			return
		}
		if argsIndexes := fastArgsIndex(buff[4:], int(numArgs)); argsIndexes != nil {
			for _, index := range argsIndexes {
				start := 4 + index[0]
				end := 4 + index[1]
				C.Args = append(C.Args, buff[start:end])
			}
		}
	}
	return C
}

func validRequest(buff []byte) bool {
	//First 4 bytes from client will always look the same
	if _, err := strconv.Atoi(string(buff[1])); err == nil && buff[0] == '*' && bytes.Equal(buff[2:4], []byte{'\r', '\n'}) {
		return true
	}
	return false
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
