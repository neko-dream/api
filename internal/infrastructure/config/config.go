package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL        string `mapstructure:"DATABASE_URL"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`
	GoogleIssuer       string `mapstructure:"GOOGLE_ISSUER"`
	DOMAIN             string `mapstructure:"DOMAIN"`
	PORT               string `mapstructure:"PORT"`
}

func LoadConfig() *Config {
	// .envファイルを読み込む
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// 環境変数を読み込む
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// .envファイルが存在しない場合は環境変数を読み込む
		default:
			panic(fmt.Errorf("設定ファイルの読み込みエラー: %w", err))
		}
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("設定ファイルの読み込みエラー: %w", err))
	}

	return &config
}
