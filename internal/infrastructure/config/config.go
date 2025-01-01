package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Env         string `mapstructure:"ENV"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`

	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleCallbackURL  string `mapstructure:"GOOGLE_CALLBACK_URL"`

	LineClientID     string `mapstructure:"LINE_CHANNEL_ID"`
	LineClientSecret string `mapstructure:"LINE_CHANNEL_SECRET"`
	LineCallbackURL  string `mapstructure:"LINE_CALLBACK_URL"`

	DOMAIN string `mapstructure:"DOMAIN"`
	PORT   string `mapstructure:"PORT"`

	TokenSecret string `mapstructure:"TOKEN_SECRET"`

	TokenPrivateKey string `mapstructure:"TOKEN_PRIVATE"`
	TokenPublicKey  string `mapstructure:"TOKEN_PUBLIC"`

	AWS_REGION            string `mapstructure:"AWS_REGION"`
	AWS_ACCESS_KEY_ID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWS_SECRET_ACCESS_KEY string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AWS_S3_ENDPOINT       string `mapstructure:"AWS_S3_ENDPOINT"`
	AWS_S3_BUCKET         string `mapstructure:"AWS_S3_BUCKET"`
	IMAGE_DOMAIN          string `mapstructure:"IMAGE_DOMAIN"`

	ANALYSIS_USER          string `mapstructure:"ANALYSIS_USER"`
	ANALYSIS_USER_PASSWORD string `mapstructure:"ANALYSIS_USER_PASSWORD"`
	ANALYSIS_API_DOMAIN    string `mapstructure:"ANALYSIS_API_DOMAIN"`

	SENTRY_DSN string `mapstructure:"SENTRY_DSN"`
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
			if err := viper.BindEnv("ENV"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("DATABASE_URL"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("GOOGLE_CLIENT_ID"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("GOOGLE_CLIENT_SECRET"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("GOOGLE_CALLBACK_URL"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("GOOGLE_ISSUER"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}

			if err := viper.BindEnv("LINE_CHANNEL_ID"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("LINE_CHANNEL_SECRET"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("LINE_CALLBACK_URL"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}

			if err := viper.BindEnv("DOMAIN"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("PORT"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}

			if err := viper.BindEnv("TOKEN_SECRET"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}

			if err := viper.BindEnv("AWS_REGION"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("AWS_ACCESS_KEY_ID"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("AWS_SECRET_ACCESS_KEY"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("AWS_S3_ENDPOINT"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("AWS_S3_BUCKET"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("IMAGE_DOMAIN"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("ANALYSIS_USER"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("ANALYSIS_USER_PASSWORD"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("ANALYSIS_API_DOMAIN"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
			if err := viper.BindEnv("SENTRY_DSN"); err != nil {
				panic(fmt.Errorf("環境変数のバインドエラー: %w", err))
			}
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
