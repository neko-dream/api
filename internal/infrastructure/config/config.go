package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/neko-dream/server/pkg/utils"
)

type Config struct {
	Env         ENV    `env:"ENV"`
	DatabaseURL string `env:"DATABASE_URL"`

	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GoogleCallbackURL  string `env:"GOOGLE_CALLBACK_URL"`

	LineClientID     string `env:"LINE_CHANNEL_ID"`
	LineClientSecret string `env:"LINE_CHANNEL_SECRET"`
	LineCallbackURL  string `env:"LINE_CALLBACK_URL"`

	DOMAIN string `env:"DOMAIN"`
	PORT   string `env:"PORT"`

	TokenSecret string `env:"TOKEN_SECRET"`

	TokenPrivateKey string `env:"TOKEN_PRIVATE"`
	TokenPublicKey  string `env:"TOKEN_PUBLIC"`

	R2_REGION            string `env:"R2_REGION"`
	R2_ACCESS_KEY_ID     string `env:"R2_ACCESS_KEY_ID"`
	R2_SECRET_ACCESS_KEY string `env:"R2_SECRET_ACCESS_KEY"`
	AWS_S3_ENDPOINT      string `env:"AWS_S3_ENDPOINT"`
	AWS_S3_BUCKET        string `env:"AWS_S3_BUCKET"`
	IMAGE_DOMAIN         string `env:"IMAGE_DOMAIN"`

	AWS_ACCESS_KEY_ID     string `env:"AWS_ACCESS_KEY_ID"`
	AWS_SECRET_ACCESS_KEY string `env:"AWS_SECRET_ACCESS_KEY"`

	ANALYSIS_USER          string `env:"ANALYSIS_USER"`
	ANALYSIS_USER_PASSWORD string `env:"ANALYSIS_USER_PASSWORD"`
	ANALYSIS_API_DOMAIN    string `env:"ANALYSIS_API_DOMAIN"`

	SENTRY_DSN       string `env:"SENTRY_DSN"`
	BASELIME_API_KEY string `env:"BASELIME_API_KEY"`

	// 暗号化バージョン (v1, etc..)
	ENCRYPTION_VERSION string `env:"ENCRYPTION_VERSION"`
	// 暗号化キー (16バイトの文字列) コマンド: openssl rand -base64 16
	ENCRYPTION_SECRET string `env:"ENCRYPTION_SECRET"`

	// ポリシーバージョン 現状は固定。API生えたら別途取得する
	POLICY_VERSION string `env:"POLICY_VERSION"`
	// ポリシー作成日 現状は固定。API生えたら別途取得する
	POLICY_CREATED_AT string `env:"POLICY_CREATED_AT"`

	EMAIL_FROM  string `env:"EMAIL_FROM"`
	APP_NAME    string `env:"APP_NAME"`
	WEBSITE_URL string `env:"WEBSITE_URL"`

	HASH_PEPPER     string `env:"HASH_PEPPER"`
	HASH_ITERATIONS int    `env:"HASH_ITERATIONS"`

	// HTTPサーバー設定
	HTTPReadTimeout  int `env:"HTTP_READ_TIMEOUT" envDefault:"15"`   // 秒
	HTTPWriteTimeout int `env:"HTTP_WRITE_TIMEOUT" envDefault:"15"`  // 秒
	HTTPIdleTimeout  int `env:"HTTP_IDLE_TIMEOUT" envDefault:"60"`   // 秒
}

type ENV string

const (
	PROD  ENV = "production"
	DEV   ENV = "development"
	LOCAL ENV = "local"
)

func (e ENV) String() string {
	return string(e)
}

func LoadConfig() *Config {
	utils.LoadEnv()

	config, err := env.ParseAs[Config]()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	return &config
}
