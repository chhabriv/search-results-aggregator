package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/chhabriv/search-results-aggregator/providerconnect"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ReqHandler interface {
	GetAggregatedSearchLinks(*gin.Context)
	CheckHealth(*gin.Context)
}

type reqHandler struct {
	providerConnect providerconnect.ProviderService
}

func NewReqHandler(providerconnect providerconnect.ProviderService) ReqHandler {
	return &reqHandler{
		providerConnect: providerconnect,
	}
}

func (h reqHandler) GetAggregatedSearchLinks(ginCtx *gin.Context) {
	log.Debug().Msg("GetAggregatedSearchLinks invoked")
	ctx := ginCtx.Request.Context()
	sortKey := strings.TrimSpace(ginCtx.Query("sortKey"))
	limit := strings.TrimSpace(ginCtx.Query("limit"))

	validationErrors := validateRequestParameters(sortKey, limit)
	if len(validationErrors) > 0 {
		log.Warn().Interface("validation_errors", validationErrors).Msg("Bad Request")
		ginCtx.JSON(http.StatusBadRequest, InvalidRequestError{Errors: validationErrors})
		return
	}

	queryProvidersResponse := h.providerConnect.QueryProviders(ctx)
	apiLinks, apiErrors := parseProviderResponse(queryProvidersResponse)
	fetchedLinksCount := len(apiLinks)
	if fetchedLinksCount == 0 && len(apiErrors) > 0 {
		ginCtx.JSON(http.StatusInternalServerError, InternalError{Errors: apiErrors})
		return
	}
	if fetchedLinksCount == 0 {
		ginCtx.JSON(http.StatusOK, GetAggregatedSearchLinksResponse{})
		return
	}

	sortLinksBySortKey(apiLinks, sortKey)
	limitInt := computeResponseLimit(fetchedLinksCount, limit)
	log.Debug().Msg("GetAggregatedSearchLinks completed")
	ginCtx.JSON(http.StatusOK, GetAggregatedSearchLinksResponse{
		Data:   apiLinks[:limitInt],
		Count:  int32(limitInt),
		Errors: apiErrors,
	})
}

func parseProviderResponse(getIndexedLinksResponse providerconnect.QueryProvidersResponse) ([]Link, []APIError) {
	apiErrors := []APIError{}
	apiLinks := []Link{}
	for _, providerResponse := range getIndexedLinksResponse.ProviderResponses {
		if providerResponse.Error != nil {
			log.Err(providerResponse.Error).Str("provider", providerResponse.Provider).Msg("failed to retrieve indexed links")
			errMsg := fmt.Sprintf("provider %s failed", providerResponse.Provider)
			apiErrors = append(apiErrors, createAPIError(APIErrorCodeProviderUnavailable, errMsg))
			continue
		}

		if len(providerResponse.Data) != 0 {
			apiLinks = append(apiLinks, convertProviderLinksToAPILinks(providerResponse.Data)...)
		}
	}
	return apiLinks, apiErrors
}

func computeResponseLimit(fetchedLinksCount int, limit string) int {
	// error is ignored as the limit is validated in validateRequestParameters
	limitInt, _ := strconv.Atoi(limit)
	if limitInt > fetchedLinksCount {
		log.Debug().Str("limit", limit).Int("links_count", fetchedLinksCount).Msg("limit is greater than the number of links, setting limit to the number of links available")
		limitInt = fetchedLinksCount
	}
	return limitInt
}

func convertProviderLinksToAPILinks(providerLinks []providerconnect.Link) []Link {
	apiLinks := []Link{}
	for _, providerLink := range providerLinks {
		apiLinks = append(apiLinks, Link{
			Url:            providerLink.Url,
			Views:          providerLink.Views,
			RelevanceScore: providerLink.RelevanceScore,
		})
	}
	return apiLinks
}
