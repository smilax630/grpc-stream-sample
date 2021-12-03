package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var (
	conf Config
)

type (
	// Config is
	Config struct {
		Redis Redis
		App   App
	}
	Redis struct {
		Host      string
		Password  string
		DB        int
		PoolSize  int
		KeyPrefix string
	}
	App struct {
		ProjectID string
	}
)

// SetupConfig is config setting
func SetupConfig() {
	if err := viper.BindEnv("env", "ENV"); err != nil {
		log.Fatalf("env bind error: %s", err)
	}
	env := viper.GetString("env")
	viper.SetConfigName(env) // 環境変数 local or dev or prod
	viper.SetConfigType("yml")
	if env == "local" {
		_, b, _, _ := runtime.Caller(0)
		Root := filepath.Join(filepath.Dir(b), "../../../")
		viper.AddConfigPath(Root + "/streamer/config/env/")
	} else {
		viper.AddConfigPath("/usr/local/bin/server/config/env/")
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("config read error: %s", err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("config error: %s", err)
	}
}

// GetRedisConfig is get redis config
func GetRedisConfig() *Redis {
	return &conf.Redis
}

// GetApp is get app config
func GetApp() *App {
	return &conf.App
}
