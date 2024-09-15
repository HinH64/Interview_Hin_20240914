package services

import (
	"Interview_Hin_20240914/models"

	"github.com/spf13/viper"
)

var Config *models.EnvConfig

func LoadConfig(configPath ...string) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("MODE", "debug")
	v.SetConfigType("dotenv")
	v.SetConfigName(".env")

	if len(configPath) > 0 {
		v.SetConfigFile(configPath[0])
	} else {
		v.AddConfigPath("./")
	}

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&Config); err != nil {
		panic(err)
	}

	if err := Config.Validate(); err != nil {
		panic(err)
	}
}
