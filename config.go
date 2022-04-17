package main

import (
	"github.com/smmr-software/mabel/internal/styles"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

type config struct {
	Download string
	Port     uint
	Theme    *styles.ColorTheme
}

type rawConfig struct {
	Download string
	Port     uint
	Theme    toml.Primitive
}

func getConfig() config {
	var conf rawConfig
	theme := &styles.DefaultTheme

	file := xdg.ConfigHome + "/mabel/config.toml"
	md, err := toml.DecodeFile(file, &conf)
	if err != nil {
		return config{Theme: theme}
	}

	if md.Type("theme") == "String" {
		var str string
		md.PrimitiveDecode(conf.Theme, &str)
		theme = styles.StringToTheme(&str)
	} else if md.Type("theme") == "Hash" {
		var hash styles.CustomTheme
		md.PrimitiveDecode(conf.Theme, &hash)
		theme = styles.CustomToTheme(&hash)
	}

	return config{
		Download: conf.Download,
		Port:     conf.Port,
		Theme:    theme,
	}
}
