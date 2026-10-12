package bootstrap

import (
	"net/http"
	"slices"
	"strings"

	"github.com/neko-dream/api/internal/infrastructure/di"
	"github.com/neko-dream/api/internal/infrastructure/http/middleware"
	"github.com/neko-dream/api/internal/presentation/handler"
	"github.com/neko-dream/api/internal/presentation/oas"
	"github.com/rs/cors"
	"go.opentelemetry.io/otel/trace"
)

// Route ルート定義を表す
type Route struct {
	Pattern     string
	Handler     http.Handler
	StripPrefix string // プレフィックスを削除する場合に指定
}

// setupRoutes ルーティングを設定する
func (b *Bootstrap) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// ルート定義
	routes := []Route{
		// API routes
		{Pattern: "/", Handler: b.setupAPIHandler()},
		// Static files
		{Pattern: "/static/", Handler: handler.NewStaticHandler(b.config), StripPrefix: "/static/"},
	}

	routes = append(routes, b.getAdminRoutes()...)
	routes = append(routes, b.getSwaggerRoutes()...)

	for _, route := range routes {
		if route.StripPrefix != "" {
			mux.Handle(route.Pattern, http.StripPrefix(route.StripPrefix, route.Handler))
		} else {
			mux.Handle(route.Pattern, route.Handler)
		}
	}

	return mux
}

// setupAPIHandler API用のハンドラーを設定して返す（CORS、認証、トレーシングなど）
func (b *Bootstrap) setupAPIHandler() http.Handler {
	srv, err := oas.NewServer(
		di.Invoke[oas.Handler](b.container),
		di.Invoke[oas.SecurityHandler](b.container),
		oas.WithErrorHandler(handler.CustomErrorHandler),
	)
	if err != nil {
		panic(err) // 初期化エラーは致命的
	}

	routeFinder := middleware.MakeRouteFinder(srv)
	tracerProvider := di.Invoke[trace.TracerProvider](b.container)

	api := middleware.Wrap(
		srv,
		middleware.Instrument("kotohiro-api", routeFinder, tracerProvider),
		middleware.SetContextCookieKey(b.container),
		middleware.Labeler(routeFinder),
	)
	return b.wrapAPIOrigins(api)
}

// wrapAPIOrigins selects the dedicated administrator guard before the ordinary API guard.
func (b *Bootstrap) wrapAPIOrigins(api http.Handler) http.Handler {
	normal := b.apiCORS().Handler(b.browserOriginGuard(api))
	adminAPI := b.adminOriginGuard(api)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/v1/admin/") || strings.HasPrefix(r.URL.Path, "/v1/manage/") {
			adminAPI.ServeHTTP(w, r)
			return
		}
		normal.ServeHTTP(w, r)
	})
}

func (b *Bootstrap) adminOriginGuard(next http.Handler) http.Handler {
	origins := b.config.AdminAllowedOrigins
	guard := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions && (origin == "" || !slices.Contains(origins, origin)) {
			http.Error(w, "Forbidden Origin", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
	return cors.New(cors.Options{AllowedOrigins: origins, AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, AllowedHeaders: []string{"Content-Type"}, AllowCredentials: true}).Handler(guard)
}

// CORS protects response reads, not simple cross-site form writes. Reject
// untrusted browser origins for mutations as well, including SameSite=None dev.
func (b *Bootstrap) browserOriginGuard(next http.Handler) http.Handler {
	origins := b.config.BrowserOrigins()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			if origin := r.Header.Get("Origin"); origin != "" && !slices.Contains(origins, origin) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`{"code":"AUTH-FORBIDDEN-ORIGIN","message":"許可されていないアクセス元です。"}`))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (b *Bootstrap) apiCORS() *cors.Cors {
	origins := b.config.BrowserOrigins()
	return cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return slices.Contains(origins, origin) },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
}

// getAdminRoutes 管理画面のルート定義を返す
func (b *Bootstrap) getAdminRoutes() []Route {
	return []Route{
		{Pattern: "/admin/", Handler: handler.NewAdminUIHandler(), StripPrefix: "/admin/"},
		{Pattern: "/admin/assets/", Handler: handler.NewAdminUIAssetsHandler(), StripPrefix: "/admin/assets/"},
		{Pattern: "/admin", Handler: http.RedirectHandler("/admin/", http.StatusSeeOther)},
	}
}

// getSwaggerRoutes Swagger UIのルート定義を返す
func (b *Bootstrap) getSwaggerRoutes() []Route {
	return []Route{
		{Pattern: "/docs/", Handler: http.HandlerFunc(b.createSwaggerHandler())},
	}
}
