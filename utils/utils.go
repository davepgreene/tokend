package utils

import (
	"reflect"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GetFunctionName uses reflection to get the name of a function as a string
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// GetTypeName uses reflection to get the name of a type as a string
func GetTypeName(i interface{}) string {
	return reflect.TypeOf(i).Name()
}

// GetLogLevel retrieves the desired log level from settings.
//
// NOTE: This should only be called after viper initializes
func GetLogLevel() log.Level {
	if lvl, err := log.ParseLevel(viper.GetString("log.level")); err == nil {
		return lvl
	}

	log.Info("Unable to parse log level. Defaulting to INFO")
	return log.InfoLevel
}

// GetLogFormatter  retrieves the desired log formatter from settings.
//
// NOTE: This should only be called after viper initializes
func GetLogFormatter() log.Formatter {
	fmt := viper.GetBool("log.json")
	if fmt {
		return &log.JSONFormatter{}
	}
	return &log.TextFormatter{}
}
