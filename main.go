package main

import (
	"github.com/pkg/profile"
	_ "net/http/pprof"
)

func main() {
	defer profile.Start(profile.ProfilePath(".")).Stop()
	s := NewServer()
	s.Init()

	s.ListenAndServerRESP()

}
