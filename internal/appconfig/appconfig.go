package appconfig

import (
	"bytes"
	"os"

	"github.com/aaabhilash97/aadhaar_scrapper_apis/configs"
	"github.com/aaabhilash97/aadhaar_scrapper_apis/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Init() *Config {
	fn := "appconfig.Init"

	log := logger.InitLogger(logger.Config{
		LogLevel:       "debug",
		LogOutputPaths: "stdout",
	})

	var config Config
	{
		viperConfig := viper.New()
		viperConfig.SetConfigType("yaml")
		err := viperConfig.ReadConfig(bytes.NewBuffer(configs.GetDefaultConfig()))
		if err != nil {
			log.Fatal(fn, zap.Error(err))
		}
		if confName := os.Getenv("conf"); confName != "" {
			envConfig, err := configs.GetYamlConfig(confName)
			if err != nil {
				log.Fatal(fn, zap.NamedError("config not found", err))
			}
			err = viperConfig.MergeConfig(
				bytes.NewBuffer(envConfig))
			if err != nil {
				log.Fatal(fn, zap.Error(err))
			}
		}

		err = viperConfig.Unmarshal(&config)
		if err != nil {
			log.Fatal(fn, zap.Error(err))
		}
	}
	return &config
}
