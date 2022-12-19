package providerconnect

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockServerConfig struct {
	httpStatus       int
	httpPath         string
	mockBodyFilePath string
}

func createMockServerConfig(
	httpStatus int,
	httpPath string,
	mockBodyFilePath string,
) mockServerConfig {
	return mockServerConfig{
		httpStatus:       httpStatus,
		httpPath:         httpPath,
		mockBodyFilePath: mockBodyFilePath,
	}
}

const (
	testGooglePath     = "/assignment132/assignment/main/google.json"
	testWikipediaPath  = "/assignment132/assignment/main/wikipedia.json"
	testDuckDuckGoPath = "/assignment132/assignment/main/duckduckgo.json"
)

func Test_QueryProviders_Success(t *testing.T) {
	mockConfigs := []mockServerConfig{
		createMockServerConfig(http.StatusOK, testGooglePath, "./testdata/google.json"),
		createMockServerConfig(http.StatusOK, testWikipediaPath, "./testdata/wikipedia.json"),
		createMockServerConfig(http.StatusOK, testDuckDuckGoPath, "./testdata/duckduckgo.json"),
	}
	mockServer := setupMockTestHTTPServer(t, mockConfigs)
	defer mockServer.Close()
	providerService := providerService{
		providerUrlMap: map[string]string{
			ProviderGoogle:     mockServer.URL + testGooglePath,
			ProviderDuckDuckGo: mockServer.URL + testDuckDuckGoPath,
			ProviderWikipedia:  mockServer.URL + testWikipediaPath,
		},
	}
	resp := providerService.QueryProviders(context.TODO())
	assert.Equal(t, 3, len(resp.ProviderResponses))
	for _, providerResp := range resp.ProviderResponses {
		if providerResp.Provider == "duckduckgo" {
			assert.Equal(t, 5, len(providerResp.Data))
			continue
		}
		if providerResp.Provider == "wikipedia" {
			assert.Equal(t, 5, len(providerResp.Data))
			continue
		}
		if providerResp.Provider == "google" {
			assert.Equal(t, 5, len(providerResp.Data))
		}
	}
}

func Test_QueryProviders_PartialSuccess(t *testing.T) {
	mockConfigs := []mockServerConfig{
		createMockServerConfig(http.StatusNotFound, testGooglePath, ""),
		createMockServerConfig(http.StatusOK, testWikipediaPath, "./testdata/wikipedia.json"),
		createMockServerConfig(http.StatusOK, testDuckDuckGoPath, "./testdata/duckduckgo.json"),
	}
	mockServer := setupMockTestHTTPServer(t, mockConfigs)
	defer mockServer.Close()
	providerService := providerService{
		providerUrlMap: map[string]string{
			ProviderGoogle:     mockServer.URL + testGooglePath,
			ProviderDuckDuckGo: mockServer.URL + testDuckDuckGoPath,
			ProviderWikipedia:  mockServer.URL + testWikipediaPath,
		},
	}
	resp := providerService.QueryProviders(context.TODO())
	assert.Equal(t, 3, len(resp.ProviderResponses))
	t.Logf("resp: %+v", resp)
	for _, providerResp := range resp.ProviderResponses {
		if providerResp.Provider == ProviderDuckDuckGo {
			assert.Equal(t, 5, len(providerResp.Data))
			continue
		}
		if providerResp.Provider == ProviderWikipedia {
			assert.Equal(t, 5, len(providerResp.Data))
			continue
		}
		if providerResp.Provider == ProviderGoogle {
			assert.Empty(t, providerResp.Data)
			expErr := "failed to retrieve links from provider: google error: unexpected end of JSON input"
			assert.EqualError(t, providerResp.Error, expErr)
		}
	}
}

func Test_QueryProviders_AllError(t *testing.T) {
	mockConfigs := []mockServerConfig{
		createMockServerConfig(http.StatusNotFound, testGooglePath, ""),
		createMockServerConfig(http.StatusOK, testWikipediaPath, ""),
		createMockServerConfig(http.StatusOK, testDuckDuckGoPath, ""),
	}
	mockServer := setupMockTestHTTPServer(t, mockConfigs)
	defer mockServer.Close()
	providerService := providerService{
		providerUrlMap: map[string]string{
			ProviderGoogle:     mockServer.URL + testGooglePath,
			ProviderDuckDuckGo: mockServer.URL + testDuckDuckGoPath,
			ProviderWikipedia:  mockServer.URL + testWikipediaPath,
		},
	}
	resp := providerService.QueryProviders(context.TODO())
	assert.Equal(t, 3, len(resp.ProviderResponses))
	t.Logf("resp: %+v", resp)
	for _, providerResp := range resp.ProviderResponses {
		if providerResp.Provider == ProviderDuckDuckGo {
			assert.Empty(t, providerResp.Data)
			expErr := "failed to retrieve links from provider: duckduckgo error: unexpected end of JSON input"
			assert.EqualError(t, providerResp.Error, expErr)
			continue
		}
		if providerResp.Provider == ProviderWikipedia {
			assert.Empty(t, providerResp.Data)
			expErr := "failed to retrieve links from provider: wikipedia error: unexpected end of JSON input"
			assert.EqualError(t, providerResp.Error, expErr)
			continue
		}
		if providerResp.Provider == ProviderGoogle {
			assert.Empty(t, providerResp.Data)
			expErr := "failed to retrieve links from provider: google error: unexpected end of JSON input"
			assert.EqualError(t, providerResp.Error, expErr)
		}
	}
}

func setupMockTestHTTPServer(
	t *testing.T,
	mockConfigs []mockServerConfig,
) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, mockConfig := range mockConfigs {
			if strings.Contains(r.RequestURI, mockConfig.httpPath) {
				w.WriteHeader(mockConfig.httpStatus)
				fileContentBytes, _ := os.ReadFile(mockConfig.mockBodyFilePath)
				// Ignoring linter check as its a test mock
				// nolint:errcheck
				w.Write(fileContentBytes)
			}
		}
	}))
}
