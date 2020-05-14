package main

import (
	_ "net/http/pprof"
)

func main() {
	//defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	s := NewServer()
	s.Init()

	s.ListenAndServerRESP()

}
