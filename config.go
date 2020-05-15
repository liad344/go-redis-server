package main

import (
	log "github.com/sirupsen/logrus"
	. "github.com/spf13/viper"
)

type ServerCfg struct {
	Addr          string
	MaxGoRoutines int
}

func (s *Server) initConfig() {
	SetConfigName("config")
	SetConfigType("toml")
	AddConfigPath("./")
	SetDefault("addr", ":8000")
	SetDefault("maxGoRoutines", 1)

	if err := ReadInConfig(); err != nil {
		log.Error("Could not load config, using default ", err)
	}
	err := Unmarshal(&s.cfg)
	if err != nil {
		log.Error("Could not load config, using default ", err)
	}
}
