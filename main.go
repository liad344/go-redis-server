package main

import (
	_ "net/http/pprof"
	"github.com/pkg/profile"
)

func main(){
	defer profile.Start(profile.ProfilePath(".")).Stop()
	s := NewServer()
	s.Init()

	s.ListenAndServerRESP()

}





