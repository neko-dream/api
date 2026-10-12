package bootstrap

import (
	"net/http"

	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v5emb"
)

// createSwaggerHandler 動的なSwaggerハンドラーを作成する
func (b *Bootstrap) createSwaggerHandler() http.HandlerFunc {
	// デフォルトのドメインを設定
	defaultDomain := b.getDefaultSwaggerDomain()

	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのHostヘッダーに基づいてドメインを決定
		domain := defaultDomain
		if b.config.Env == config.DEV {
			host := r.Host
			if host == "api-dev.kotohiro.com" || host == "api.dev.kotohiro.com" {
				domain = "https://" + host + "/static/oas/openapi.yaml"
			}
		}

		swagger := b.createSwaggerUI(domain)
		swagger.ServeHTTP(w, r)
	}
}

// getDefaultSwaggerDomain デフォルトのSwaggerドメインを返す
func (b *Bootstrap) getDefaultSwaggerDomain() string {
	switch b.config.Env {
	case config.DEV:
		return "https://api-dev.kotohiro.com/static/oas/openapi.yaml"
	case config.PROD:
		return "https://api.kotohiro.com/static/oas/openapi.yaml"
	default:
		return "http://localhost:" + b.config.PORT + "/static/oas/openapi.yaml"
	}
}

// createSwaggerUI Swagger UIの設定を作成する
func (b *Bootstrap) createSwaggerUI(domain string) http.Handler {
	tagsSorterFunc := "(a, b) => {" +
		"const priority = {\"auth\": 1, \"user\": 2, \"talk_session\": 3, \"opinion\": 4, \"organization\": 5, \"vote\": 6}; " +
		"const ap = priority[a.toLowerCase()]; " +
		"const bp = priority[b.toLowerCase()]; " +
		"if (ap && bp) return ap - bp; " +
		"if (ap) return -1; " +
		"if (bp) return 1; " +
		"return a.toLowerCase().localeCompare(b.toLowerCase());" +
		"}"

	swagger := v5emb.NewWithConfig(swgui.Config{
		Title:       "Kotohiro API",
		HideCurl:    true,
		SwaggerJSON: domain,
		BasePath:    "/docs/",
		JsonEditor:  true,
		SettingsUI: map[string]string{
			"deepLinking":              "true",
			"defaultModelsExpandDepth": "-1",
			"defaultModelExpandDepth":  "-1",
			"defaultModelRendering":    "\"model\"",
			"displayRequestDuration":   "true",
			"tryItOutEnabled":          "true",
			"layout":                   "\"BaseLayout\"",
			"showExtensions":           "true",
			"showCommonExtensions":     "true",
			"syntaxHighlight":          "{\"activate\": true,\"theme\": \"tomorrow-night\"}",
			"displayOperationId":       "true",
			"filter":                   "true",
			"operationsSorter":         "\"alpha\"",
			"tagsSorter":               tagsSorterFunc,
		},
	})

	return swagger("Kotohiro API", domain, "/docs/")
}
