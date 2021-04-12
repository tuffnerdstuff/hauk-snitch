package config

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/tuffnerdstuff/hauk-snitch/hauk"
	"github.com/tuffnerdstuff/hauk-snitch/mqtt"
)

// LoadConfig loads config.toml
func LoadConfig() {

	setMqttDefaults()
	setHaukDefaults()
	readConfigFromFile()
}

// GetMqttConfig returns a struct containing mqtt config values
func GetMqttConfig() mqtt.Config {
	var mqttConfig mqtt.Config
	mqttConfig.Host = viper.GetString("mqtt.host")
	mqttConfig.Port = viper.GetInt("mqtt.port")
	mqttConfig.User = viper.GetString("mqtt.user")
	mqttConfig.Password = viper.GetString("mqtt.password")
	mqttConfig.IsAnonymous = viper.GetBool("mqtt.anonymous")
	mqttConfig.IsTLS = viper.GetBool("mqtt.tls")
	return mqttConfig
}

// GetHaukConfig returns a struct containing hauk config values
func GetHaukConfig() hauk.Config {
	var haukConfig hauk.Config
	haukConfig.Host = viper.GetString("hauk.host")
	haukConfig.Port = viper.GetInt("hauk.port")
	haukConfig.User = viper.GetString("hauk.user")
	haukConfig.Password = viper.GetString("hauk.password")
	haukConfig.IsAnonymous = viper.GetBool("hauk.anonymous")
	haukConfig.IsTLS = viper.GetBool("hauk.tls")
	haukConfig.Duration = viper.GetInt("hauk.duration")
	haukConfig.Interval = viper.GetInt("hauk.interval")
	return haukConfig
}

func readConfigFromFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Config error: %w", err))
	}
}

func setMqttDefaults() {
	viper.SetDefault("mqtt.host", "localhost")
	viper.SetDefault("mqtt.port", 1883)
	viper.SetDefault("mqtt.user", "")
	viper.SetDefault("mqtt.password", "")
	viper.SetDefault("mqtt.anonymous", true)
	viper.SetDefault("mqtt.tls", false)
}

func setHaukDefaults() {
	viper.SetDefault("hauk.host", "localhost")
	viper.SetDefault("hauk.port", 80)
	viper.SetDefault("hauk.user", "")
	viper.SetDefault("hauk.password", "")
	viper.SetDefault("hauk.anonymous", true)
	viper.SetDefault("hauk.tls", false)
	viper.SetDefault("hauk.duration", 3600) // 1 hour
	viper.SetDefault("hauk.interval", 1)    // Every second
}
