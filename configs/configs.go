package configs

import (
	_ "embed"
)

//go:embed default.yml
var defaultConf []byte

func GetDefaultConfig() []byte {
	return defaultConf
}
