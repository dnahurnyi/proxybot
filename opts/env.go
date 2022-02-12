package opts

import (
	"fmt"

	"github.com/spf13/viper"
)

func ReadOS() (*Config, error) {
	err := Load()
	if err != nil {
		return nil, fmt.Errorf("load configs: %w", err)
	}

	viper.AutomaticEnv()

	viper.SetEnvPrefix("APP")
	viper.SetDefault("APP_NAME", "tg-proxy")

	return &Config{
		AppName:      viper.GetString("APP_NAME"),
		AppID:        viper.GetInt32("TG_ID"),
		MasterChatID: viper.GetInt64("MASTER_CHAT_ID"),
		AppHash:      viper.GetString("TG_HASH"),
		DB: DB{
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			DBName:   viper.GetString("DB_NAME"),
			SSL:      viper.GetBool("DB_SSL"),
		},
	}, nil
}
