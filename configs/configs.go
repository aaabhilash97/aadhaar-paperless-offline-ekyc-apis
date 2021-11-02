package configs

import (
	"embed"
	_ "embed"
	"fmt"
)

//go:embed *.yml
var f embed.FS

func GetYamlConfig(name string) ([]byte, error) {
	return f.ReadFile(fmt.Sprintf("%s.yml", name))
}

//go:embed default.yml
var defaultConf []byte

func GetDefaultConfig() []byte {
	return defaultConf
}
