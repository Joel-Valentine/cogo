package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

// Provider defines a set of read-only methods for accessing the application
// configuration params as defined in one of the config files.
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
}

// AppError is The default config error
type AppError struct {
	Error   error
	Message string
	Code    int
}

var defaultConfig *viper.Viper

// Config returns a default config providers
func Config() (Provider, *AppError) {
	if err := defaultConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, &AppError{err, "No config file", 01}
		} else {
			return nil, &AppError{err, "Unable to read config file", 02}
		}
	}
	return defaultConfig, nil
}

// LoadConfigProvider returns a configured viper instance
func LoadConfigProvider(appName string) Provider {
	return readViperConfig()
}

func init() {
	defaultConfig = readViperConfig()
}

func readViperConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName(".cogo_config")
	v.SetConfigType("json")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	// global defaults

	return v
}

// SaveConfigFile will save a config file to be used again later
func SaveConfigFile(token string) {
	defaultConfig.Set("digitalOceanToken", token)
	if err := defaultConfig.WriteConfigAs(".cogo_config.json"); err != nil {
		fmt.Printf("UH OH %s", err)
		return
	}
}
