package appconfig

import (
	"os"

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
		viperConfig.SetConfigType("yml")

		viperConfig.SetConfigName("default")
		viperConfig.AddConfigPath("./configs")
		err := viperConfig.ReadInConfig()
		if err != nil {
			log.Fatal(fn, zap.Error(err))
		}
		if context := os.Getenv("env"); context != "" {
			viperConfig.SetConfigName(context)
			err = viperConfig.MergeInConfig()
			if err != nil {
				log.Fatal(fn, zap.String("context", context), zap.Error(err))
			}
		}

		err = viperConfig.Unmarshal(&config)
		if err != nil {
			log.Fatal(fn, zap.Error(err))
		}
	}
	return &config
}
