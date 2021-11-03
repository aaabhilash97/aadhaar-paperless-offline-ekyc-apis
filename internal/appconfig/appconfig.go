package appconfig

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aaabhilash97/aadhaar-paperless-offline-ekyc-apis/configs"
	"github.com/aaabhilash97/aadhaar-paperless-offline-ekyc-apis/pkg/logger"
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
		version := pflag.Bool("version", false, "Show version information")

		pflag.String("conf", "", "config path")
		pflag.Parse()

		if *version {
			msg := fmt.Sprintf(`Git Tag:      %s
Git commit:   %s`,
				os.Getenv("VERSION_INFO_GIT_TAG"),
				os.Getenv("VERSION_INFO_GIT_COMMIT"))

			fmt.Println(msg)
			os.Exit(0)
		}

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
