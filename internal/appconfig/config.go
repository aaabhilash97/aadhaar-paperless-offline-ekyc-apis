package appconfig

import (
	"github.com/aaabhilash97/aadhaar_scrapper_apis/internal/aadhaarcache"
	"github.com/aaabhilash97/aadhaar_scrapper_apis/pkg/logger"
	"github.com/aaabhilash97/aadhaar_scrapper_apis/pkg/service/v1"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Config struct {
	Server struct {
		Env      string
		Host     string
		GrpcPort string
		HttpPort string
	}

	Logger struct {
		OutputPath        string
		Level             string
		DisableStackTrace bool
	}

	AccessLogger struct {
		OutputPath        string
		Level             string
		DisableStackTrace bool
	}

	AadhaarCacheStore struct {
		Type  string
		Redis struct {
			Host     string
			Password string
			DB       int
		}
	}

	logger       *zap.Logger
	accessLogger *zap.Logger
}

func (c Config) Close() {}

func (c Config) GetLogger() *zap.Logger {
	if c.logger != nil {
		return c.logger
	}
	c.logger = logger.InitLogger(logger.Config{
		LogLevel:          c.Logger.Level,
		LogOutputPaths:    c.Logger.OutputPath,
		DisableStackTrace: c.Logger.DisableStackTrace,
	})

	return c.logger
}

func (c Config) GetAccessLogger() *zap.Logger {
	if c.logger != nil {
		return c.accessLogger
	}
	c.accessLogger = logger.InitLogger(logger.Config{
		LogLevel:          c.AccessLogger.Level,
		LogOutputPaths:    c.AccessLogger.OutputPath,
		DisableStackTrace: c.AccessLogger.DisableStackTrace,
	})

	return c.accessLogger
}

func (c Config) GetAadhaarCacheStore() service.AadhaarCacheStore {
	if c.AadhaarCacheStore.Type == "redis" {
		return aadhaarcache.NewRedisCache(&redis.Options{
			Addr:     c.AadhaarCacheStore.Redis.Host,
			Password: c.AadhaarCacheStore.Redis.Password, // no password set
			DB:       c.AadhaarCacheStore.Redis.DB,       // use default DB
		})
	}
	return aadhaarcache.NewMemCache()
}
