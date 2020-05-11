package main

import (
	"testing"
)

func parserTest (t *testing.B){
	b := []byte("*5\\r\\n$3\\r\\nset\\r\\n$2\\r\\nyo\\r\\n$2\\r\\nyo\\r\\n$2\\r\\nex\\r\\n$1\\r\\n1\\r\\n\\")
	cmd := parseCmd(b)
	res := []string{"set" , "yo" , "yo" , "ex" , "1"}

	for i , arg := range cmd.Args{
		if string(arg) != res[i]{
			t.Fail()
		}
	}
}