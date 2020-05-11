package main

import (
	"testing"
)
func TestGet (t *testing.T) {
	var b = []byte{42, 50, 13, 10, 36, 51, 13, 10, 103, 101, 116, 13, 10, 36, 50, 13, 10, 121, 111, 13, 10}
	cmd := parseCmd(b)
	res := []string{"get", "yo"}

	for i, arg := range cmd.Args {
		if string(arg) != res[i] {
			t.Fail()
		}
	}
}