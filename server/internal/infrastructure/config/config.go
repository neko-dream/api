package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/neko-dream/api/pkg/utils"
	"net/url"
	"strings"
)

type Config struct {
	Env         ENV    `env:"ENV"`
	DatabaseURL string `env:"DATABASE_URL"`

	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GoogleCallbackURL  string `env:"GOOGLE_CALLBACK_URL"`

	// 初期管理者。起動時に admin_users に存在しなければ作成する。
	InitialAdminEmail   string   `env:"INITIAL_ADMIN_EMAIL"`
	InitialAdminName    string   `env:"INITIAL_ADMIN_NAME"`
	AdminJWTSecret      string   `env:"ADMIN_JWT_SECRET"`
	AdminAllowedOrigins []string `env:"ADMIN_ALLOWED_ORIGINS" envSeparator:","`

	// AdditionalBrowserOrigins WEBSITE_URL 以外に、ブラウザからの認証付きアクセスを許可するオリジン（例: ステージングのWeb）。
	AdditionalBrowserOrigins []string `env:"ADDITIONAL_BROWSER_ORIGINS" envSeparator:","`

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

	PINPOINT_APPLICATION_ID string `env:"PINPOINT_APPLICATION_ID"`

	// FCM VAPID Key
	FCM_VAPID_KEY string `env:"FCM_VAPID_KEY"`

	// ポリシーバージョン 現状は固定。API生えたら別途取得する
	POLICY_VERSION string `env:"POLICY_VERSION"`
	// ポリシー作成日 現状は固定。API生えたら別途取得する
	POLICY_CREATED_AT string `env:"POLICY_CREATED_AT"`

	EMAIL_FROM                  string `env:"EMAIL_FROM"`
	EMAIL_VERIFICATION_BASE_URL string `env:"EMAIL_VERIFICATION_BASE_URL"`
	EMAIL_VERIFICATION_PAGE_URL string `env:"EMAIL_VERIFICATION_PAGE_URL"`
	SMTP_ADDR                   string `env:"SMTP_ADDR" envDefault:"127.0.0.1:1025"`
	APP_NAME                    string `env:"APP_NAME"`
	WEBSITE_URL                 string `env:"WEBSITE_URL"`

	HASH_PEPPER     string `env:"HASH_PEPPER"`
	HASH_ITERATIONS int    `env:"HASH_ITERATIONS"`

	// HTTPサーバー設定
	HTTPReadTimeout  int `env:"HTTP_READ_TIMEOUT" envDefault:"15"`  // 秒
	HTTPWriteTimeout int `env:"HTTP_WRITE_TIMEOUT" envDefault:"15"` // 秒
	HTTPIdleTimeout  int `env:"HTTP_IDLE_TIMEOUT" envDefault:"60"`  // 秒
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

// Development-only authentication is opt-in, never "anything but production".
func (e ENV) AllowsDevelopmentAuthentication() bool {
	return e == DEV || e == LOCAL
}

func (e ENV) Validate() error {
	switch e {
	case PROD, DEV, LOCAL:
		return nil
	default:
		return fmt.Errorf("ENV must be production, development or local; got %q", e)
	}
}

// BrowserOrigins is shared by credentialed CORS and email-login CSRF checks.
// No arbitrary HTTPS hosts, wildcard previews, or unknown environments are trusted.
func (c *Config) BrowserOrigins() []string {
	if c == nil || c.Env.Validate() != nil {
		return nil
	}
	var origins []string
	for _, base := range []string{c.WEBSITE_URL, c.EMAIL_VERIFICATION_BASE_URL} {
		u, err := url.Parse(base)
		if err == nil && u.Scheme == "https" && u.Host != "" && u.User == nil && u.RawQuery == "" && u.Fragment == "" {
			origins = append(origins, u.Scheme+"://"+u.Host)
		}
	}
	// Both production Web domains are routed to the same app, not redirected.
	if c.Env == PROD {
		for _, origin := range origins {
			if origin == "https://www.kotohiro.com" {
				origins = append(origins, "https://kotohiro.com")
				break
			}
		}
	}
	// Legacy development Swagger/admin URLs still use this explicit API alias.
	if c.Env == DEV {
		origins = append(origins, "https://api-dev.kotohiro.com")
	}
	for _, raw := range c.AdditionalBrowserOrigins {
		if origin, ok := explicitHTTPSOrigin(raw); ok {
			origins = append(origins, origin)
		}
	}
	return origins
}

// explicitHTTPSOrigin はパス・クエリ・認証情報・ワイルドカードを含まない HTTPS オリジンだけを受け付ける。
func explicitHTTPSOrigin(raw string) (string, bool) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Scheme != "https" || u.Host == "" || strings.Contains(u.Host, "*") || u.User != nil ||
		(u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" {
		return "", false
	}
	return u.Scheme + "://" + u.Host, true
}

func LoadConfig() *Config {
	utils.LoadEnv()

	config, err := env.ParseAs[Config]()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}
	if err := config.Env.Validate(); err != nil {
		panic(err)
	}

	return &config
}
