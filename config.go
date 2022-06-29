package main

import (
	"github.com/smmr-software/mabel/full"
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

	Log               bool
	RequireEncryption bool `toml:"require_encryption"`

	Theme toml.Primitive
	Keys  full.CustomKeyMap
}

// getTheme checks the config file for a configured theme key and
// returns a ColorTheme. If the theme key is not present, the function
// returns the default theme.
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

// getConfig checks for a TOML file in the config directory and returns
// a config.
func getConfig() (conf config) {
	file := xdg.ConfigHome + "/mabel/config.toml"

	md, err = toml.DecodeFile(file, &conf)
	if err != nil {
		return config{}
	}

	return conf
}
