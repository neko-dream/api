package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/neko-dream/server/internal/infrastructure/di"
	"github.com/neko-dream/server/internal/infrastructure/middleware"
	"github.com/neko-dream/server/internal/presentation/oas"
)

func main() {
	port := os.Getenv("PORT")

	container := di.BuildContainer()
	srv, err := oas.NewServer(
		di.Invoke[oas.Handler](container),
		di.Invoke[oas.SecurityHandler](container),
	)
	if err != nil {
		panic(err)
	}

	corsHandler := middleware.CORSMiddleware(srv)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), corsHandler); err != nil {
		panic(err)
	}
}
