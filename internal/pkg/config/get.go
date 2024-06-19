package config

import (
	"time"

	"github.com/DmitySH/tg-gpt/internal/pkg/loggy"
	"github.com/spf13/viper"
)

func String(key string) string {
	val := viper.GetString(key)
	if val == "" {
		loggy.Warnln("config value for key", key, "is empty")
	}

	return val
}

func Int(key string) int {
	return viper.GetInt(key)
}

func Duration(key string) time.Duration {
	return viper.GetDuration(key)
}
