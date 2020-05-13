package main

import (
	_ "net/http/pprof"
)

func main() {
	s := NewServer()
	s.Init()

	s.ListenAndServerRESP()

}
