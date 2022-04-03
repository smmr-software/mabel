package main

import (
	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

type config struct {
	Download string
	Port     uint
}

func getConfig() config {
	var conf config

	file := xdg.ConfigHome + "/mabel/config.toml"
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		return config{}
	}

	return conf
}
