package appconfig

import (
	"bytes"
	"os"

	"github.com/aaabhilash97/aadhaar_scrapper_apis/configs"
	"github.com/aaabhilash97/aadhaar_scrapper_apis/pkg/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Init() *Config {
	fn := "appconfig.Init"

	log := logger.InitLogger(logger.Config{
		LogLevel:          "debug",
		LogOutputPaths:    "stdout",
		DisableStackTrace: true,
	})

	var config Config
	{
		pflag.String("conf", "", "config path")
		pflag.Parse()

		viperConfig := viper.New()
		viperConfig.AutomaticEnv()
		_ = viperConfig.BindPFlags(pflag.CommandLine)
		viperConfig.SetConfigType("yml")
		err := viperConfig.ReadConfig(bytes.NewBuffer(configs.GetDefaultConfig()))
		if err != nil {
			log.Fatal(fn, zap.Error(err))
		}
		if confPath := viperConfig.GetString("conf"); confPath != "" {
			confPath, err := os.Open(confPath)
			if err != nil {
				log.Fatal(fn, zap.NamedError("config not found", err))
			}
			err = viperConfig.MergeConfig(confPath)
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
