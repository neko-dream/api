package bootstrap

import (
	"net/http"

	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/http/middleware"
	"github.com/neko-dream/server/internal/presentation/handler"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/rs/cors"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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
		{Pattern: "/static/", Handler: handler.NewStaticHandler(), StripPrefix: "/static/"},
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

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	routeFinder := middleware.MakeRouteFinder(srv)
	tracerProvider := di.Invoke[*sdktrace.TracerProvider](b.container)

	return c.Handler(middleware.Wrap(
		srv,
		middleware.Instrument("kotohiro-api", routeFinder, tracerProvider),
		middleware.SetContextCookieKey(b.container),
		middleware.Labeler(routeFinder),
	))
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
