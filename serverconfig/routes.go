package serverconfig

import (
	"github.com/gin-gonic/gin"

	"github.com/chhabriv/search-results-aggregator/api"
)

func setupRoutes(router *gin.Engine, handler api.ReqHandler) {
	router.GET("/health", handler.CheckHealth)
	router.GET("/links", handler.GetAggregatedSearchLinks)
}
