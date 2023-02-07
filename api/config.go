package main

import "github.com/spf13/viper"

type Config struct {
	Token          string `mapstructure:"TELEGRAM_TOKEN"`
	DatabaseSource string `mapstructure:"DBSOURCE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
