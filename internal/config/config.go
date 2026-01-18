package config

import (
	"github.com/spf13/viper"
)

// Get returns a configuration value by key
func Get(key string) interface{} {
	return viper.Get(key)
}

// GetString returns a string configuration value by key
func GetString(key string) string {
	return viper.GetString(key)
}

// GetBool returns a bool configuration value by key
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetInt returns an int configuration value by key
func GetInt(key string) int {
	return viper.GetInt(key)
}

// Set sets a configuration value
func Set(key string, value interface{}) {
	viper.Set(key, value)
}
