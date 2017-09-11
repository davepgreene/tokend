package utils

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Defaults generates a set of default configuration options
func Defaults() {
	localhost := "127.0.0.1"

	viper.SetDefault("vault", map[string]interface{}{
		"host":                  localhost,
		"port":                  8200,
		"token_renew_increment": 60,
		"tls": false,
	})

	viper.SetDefault("metadata", map[string]interface{}{
		"host": "169.254.169.254",
		"port": 80,
	})

	viper.SetDefault("service", map[string]interface{}{
		"host":    localhost,
		"port":    4500,
		"backend": "warden",
	})

	viper.SetDefault("log", map[string]interface{}{
		"level":    log.InfoLevel,
		"json":     true,
		"requests": true,
	})

	viper.SetDefault("warden", map[string]interface{}{
		"host": localhost,
		"port": 8705,
		"path": "/v1/authenticate",
	})

	viper.SetDefault("storage", map[string]interface{}{
		"timeout": 500,
	})
}
