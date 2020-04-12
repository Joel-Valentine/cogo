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

// Configurations is the structure of the config file
type Configurations struct {
	digitalOceanAPIToken string
}

// PossibleSaveLocations is a list of all locations that is currently supported
// Not entirely sure this is what I want.. I think I want to use an enum
var PossibleSaveLocations = []string{"$HOME/.cogo", "$HOME/.config/.cogo", "./.cogo"}

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
	fmt.Printf("%s", PossibleSaveLocations)
	v := viper.New()
	v.SetConfigName(".cogo")
	v.SetConfigType("json")
	v.AddConfigPath("$HOME")
	v.AddConfigPath("$HOME/.config/")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	// global defaults

	return v
}
