package config

import (
	"strings"

	_ "github.com/joho/godotenv/autoload" // secret variable replacement for local test
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Gin struct {
		Mode string // enum: debug, release, test
	}
	Server struct {
		Port               string
		ShutdownTimeoutSec int // 0 would shut down immediately
	}
	DB struct {
		Postgres struct {
			Host     string `envconfig:"PG_HOST"`
			Port     int    `envconfig:"PG_PORT"`
			DBName   string `envconfig:"PG_DB_NAME"`
			User     string `envconfig:"PG_USER"`
			Password string `envconfig:"PG_PASSWORD"`
			SSLMode  string `envconfig:"PG_SSL_MODE"`
		}
	}
}

func New() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("Error reading env file:", zap.Error(err))
	}

	cfg := new(Config)
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Fatal("Failed to unmarshal config:", zap.Error(err))
	}

	return cfg
}
