package main

import (
	"github.com/smmr-software/mabel/internal/styles"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

var (
	md  toml.MetaData
	err error
)

type config struct {
	Download string
	Port     uint
	Log      bool
	Theme    toml.Primitive
}

// getTheme checks the config file for a configured theme and returns
// a ColorTheme from a string referring to a predefined theme or a map
// of a custom theme.
func (c *config) getTheme() *styles.ColorTheme {
	switch md.Type("theme") {
	case "String":
		var str string
		md.PrimitiveDecode(c.Theme, &str)
		return styles.StringToTheme(&str)
	case "Hash":
		var hash styles.CustomTheme
		md.PrimitiveDecode(c.Theme, &hash)
		return hash.ToTheme()
	default:
		return &styles.DefaultTheme
	}
}

// getConfig checks for a config file in the directory, reads it, and
// decodes the file to be returned.
func getConfig() (conf config) {
	file := xdg.ConfigHome + "/mabel/config.toml"

	md, err = toml.DecodeFile(file, &conf)
	if err != nil {
		return config{}
	}

	return conf
}
