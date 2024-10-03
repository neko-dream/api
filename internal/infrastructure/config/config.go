package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL        string `mapstructure:"DATABASE_URL"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleCallbackURL  string `mapstructure:"GOOGLE_CALLBACK_URL"`
	GoogleIssuer       string `mapstructure:"GOOGLE_ISSUER"`
	DOMAIN             string `mapstructure:"DOMAIN"`
	PORT               string `mapstructure:"PORT"`

	AWS_REGION            string `mapstructure:"AWS_REGION"`
	AWS_ACCESS_KEY_ID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWS_SECRET_ACCESS_KEY string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AWS_S3_ENDPOINT       string `mapstructure:"AWS_S3_ENDPOINT"`
	AWS_S3_BUCKET         string `mapstructure:"AWS_S3_BUCKET"`
	IMAGE_DOMAIN          string `mapstructure:"IMAGE_DOMAIN"`
}

func LoadConfig() *Config {
	// .envファイルを読み込む
	viper.AddConfigPath(".")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../../../")
	viper.AddConfigPath("../../../../")
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
