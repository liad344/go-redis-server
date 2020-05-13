package main

import (
	_ "net/http/pprof"
	"github.com/pkg/profile"
)

func main(){
	defer profile.Start(profile.CPUProfile).Stop()
	s := NewServer()
	s.Init()

	s.ListenAndServerRESP()

}





