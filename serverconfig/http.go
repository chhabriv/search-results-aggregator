package serverconfig

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/chhabriv/search-results-aggregator/api"
	"github.com/chhabriv/search-results-aggregator/providerconnect"
)

// TODO: Consider moving this to a config file
var providerURLMap = map[string]string{
	providerconnect.ProviderDuckDuckGo: "https://raw.githubusercontent.com/assignment132/assignment/main/duckduckgo.json",
	providerconnect.ProviderWikipedia:  "https://raw.githubusercontent.com/assignment132/assignment/main/wikipedia.json",
	providerconnect.ProviderGoogle:     "https://raw.githubusercontent.com/assignment132/assignment/main/google.json",
}

const (
	serverPort                    = ":8080"
	serverHTTPReadTimeout         = 5 * time.Second
	serverHTTPWriteTimeout        = 5 * time.Second
	serverGracefulShutdownTimeout = 5 * time.Second

	httpClientTimeout = 1 * time.Second
	httpClientRetries = 1
)

// StartServer initilizes the dependencies and starts
// http server.
func StartServer(ctx context.Context) {
	providerService := initProviderConnect()
	reqHandler := api.NewReqHandler(providerService)
	router := gin.Default()
	setupRoutes(router, reqHandler)

	srv := &http.Server{
		Addr:         serverPort,
		Handler:      router,
		ReadTimeout:  serverHTTPReadTimeout,
		WriteTimeout: serverHTTPWriteTimeout,
	}

	go func() {
		if startErr := srv.ListenAndServe(); startErr != nil && !errors.Is(startErr, http.ErrServerClosed) {
			log.Panic().Err(startErr).Msg("Server failed to start")
		}
	}()

	log.Info().Str("port", serverPort).Msg("Server started")

	<-ctx.Done()

	log.Info().Msg("Server shutting down")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), serverGracefulShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}
}

func initProviderConnect() providerconnect.ProviderService {
	return providerconnect.NewProviderService(httpClientTimeout, providerURLMap, httpClientRetries)
}
