package providerconnect

//go:generate mockgen -destination=testdata/mocks/mock_providers.go -package=mocks . ProviderService

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	ProviderDuckDuckGo = "duckduckgo"
	ProviderGoogle     = "google"
	ProviderWikipedia  = "wikipedia"
)

type providerService struct {
	httpClient     http.Client
	providerUrlMap map[string]string
	retries        int
}

type ProviderService interface {
	QueryProviders(context.Context) QueryProvidersResponse
}

func NewProviderService(httpTimeout time.Duration, providerUrlMap map[string]string, retries int) ProviderService {
	httpClient := http.Client{
		Timeout: httpTimeout,
	}
	return &providerService{
		httpClient:     httpClient,
		providerUrlMap: providerUrlMap,
		retries:        retries,
	}
}

func (ps providerService) QueryProviders(ctx context.Context) QueryProvidersResponse {
	chOut := make(chan ProviderResponse)
	var wg sync.WaitGroup

	for provider := range ps.providerUrlMap {
		wg.Add(1)
		go ps.callProvider(ctx, provider, chOut, &wg)
	}

	go func() {
		wg.Wait()
		close(chOut)
	}()
	return processProviderResponses(ctx, chOut)
}

func (ps providerService) callProvider(
	ctx context.Context,
	provider string,
	chOut chan<- ProviderResponse,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	var resp ProviderResponse
	url := ps.providerUrlMap[provider]
	clientErr := doHTTPGet(ctx, ps.httpClient, url, &resp, 1)
	if clientErr != nil {
		resp.Error = fmt.Errorf("failed to retrieve links from provider: %s error: %w", provider, clientErr)
	}
	resp.Provider = provider
	chOut <- resp
}

func processProviderResponses(
	ctx context.Context,
	chOut <-chan ProviderResponse,
) QueryProvidersResponse {
	providerResponses := []ProviderResponse{}
	for resp := range chOut {
		providerResponses = append(providerResponses, resp)
	}
	return QueryProvidersResponse{
		ProviderResponses: providerResponses,
	}
}
