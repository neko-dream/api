package bootstrap

import (
	"fmt"
	"net/http"
	"time"
)

// setupHTTPServer HTTPサーバーを設定して返す
func (b *Bootstrap) setupHTTPServer() (*http.Server, error) {
	mux := b.setupRoutes()

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", b.config.PORT),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}, nil
}

// startHTTPServer HTTPサーバーを起動する
func (b *Bootstrap) startHTTPServer() error {
	server, err := b.setupHTTPServer()
	if err != nil {
		return fmt.Errorf("failed to setup HTTP server: %w", err)
	}

	// グレースフルシャットダウンの設定は次のステップで

	return server.ListenAndServe()
}
