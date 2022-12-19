package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/chhabriv/search-results-aggregator/providerconnect"
	"github.com/chhabriv/search-results-aggregator/providerconnect/testdata/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testSetup struct {
	recorder     *httptest.ResponseRecorder
	ginRouter    *gin.Engine
	ginContext   *gin.Context
	mockProvider *mocks.MockProviderService
	goMockCtrl   *gomock.Controller
	reqHandler   ReqHandler
}

func setupTest(t *testing.T) testSetup {
	var ts testSetup
	ts.recorder = httptest.NewRecorder()
	ts.ginContext, ts.ginRouter = gin.CreateTestContext(ts.recorder)
	ts.goMockCtrl = gomock.NewController(t)
	ts.mockProvider = mocks.NewMockProviderService(ts.goMockCtrl)
	ts.reqHandler = NewReqHandler(ts.mockProvider)
	ts.ginRouter.GET("/links", ts.reqHandler.GetAggregatedSearchLinks)
	return ts
}

func createValidRequest(sortBy, limit string) (*http.Request, error) {
	request, httpRequestErr := http.NewRequest("GET", "/links", nil)
	u := url.Values{}
	u.Add("sortKey", sortBy)
	u.Add("limit", limit)
	request.URL.RawQuery = u.Encode()
	return request, httpRequestErr
}

func Test_GetAggregatedSearchLinks_Validation_Error(t *testing.T) {
	// test setup
	ts := setupTest(t)
	defer ts.goMockCtrl.Finish()

	// create http request
	request, errHttp := http.NewRequest("GET", "/links", nil)
	assert.NoError(t, errHttp)

	// serve http request
	ts.ginRouter.ServeHTTP(ts.recorder, request)

	// assert response
	expHTTPStatus := http.StatusBadRequest
	assert.Equal(t, expHTTPStatus, ts.recorder.Code)

	actualResponse := InvalidRequestError{}
	errUnmarshal := json.Unmarshal(ts.recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, errUnmarshal)
	expBadRequestErrors := InvalidRequestError{
		Errors: []BadRequestError{
			{
				Field:       "sortKey",
				Description: "sortKey is invalid. It should be either relevanceScore or views",
			},
			{
				Field:       "limit",
				Description: "limit should be a number greater than 1 and less than 200",
			},
		},
	}
	assert.Equal(t, expBadRequestErrors, actualResponse)
}

func Test_GetAggregatedSearchLinks_ProviderConnect_EmptyResponse(t *testing.T) {
	ts := setupTest(t)
	defer ts.goMockCtrl.Finish()

	request, errHttp := createValidRequest(sortKeyRelevanceScore, "3")
	assert.NoError(t, errHttp)

	ts.mockProvider.EXPECT().QueryProviders(gomock.Any()).Return(providerconnect.QueryProvidersResponse{
		ProviderResponses: []providerconnect.ProviderResponse{},
	})

	ts.ginRouter.ServeHTTP(ts.recorder, request)

	var actualResponse GetAggregatedSearchLinksResponse
	errUnmarshal := json.Unmarshal(ts.recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, errUnmarshal)

	expHTTPStatus := http.StatusOK
	assert.Equal(t, expHTTPStatus, ts.recorder.Code)

	expResponse := GetAggregatedSearchLinksResponse{}
	assert.Equal(t, expResponse, actualResponse)
}

func Test_GetAggregatedSearchLinks_ProviderConnect_AllError(t *testing.T) {
	ts := setupTest(t)
	defer ts.goMockCtrl.Finish()

	request, errHttp := createValidRequest(sortKeyRelevanceScore, "3")
	assert.NoError(t, errHttp)

	ts.mockProvider.EXPECT().QueryProviders(gomock.Any()).Return(mockProviderConnectErrorResponse())

	ts.ginRouter.ServeHTTP(ts.recorder, request)

	var actualResponse InternalError
	errUnmarshal := json.Unmarshal(ts.recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, errUnmarshal)

	expHTTPStatus := http.StatusInternalServerError
	assert.Equal(t, expHTTPStatus, ts.recorder.Code)

	expErrResponse := InternalError{
		Errors: []APIError{
			{
				Code:    APIErrorCodeProviderUnavailable,
				Message: "provider duckduckgo failed",
			},
			{
				Code:    APIErrorCodeProviderUnavailable,
				Message: "provider google failed",
			},
			{
				Code:    APIErrorCodeProviderUnavailable,
				Message: "provider wikipedia failed",
			},
		},
	}
	assert.Equal(t, expErrResponse, actualResponse)
}

func Test_GetAggregatedSearchLinks_ProviderConnect_PartialError(t *testing.T) {
	ts := setupTest(t)
	defer ts.goMockCtrl.Finish()

	request, errHttp := createValidRequest(sortKeyRelevanceScore, "3")
	assert.NoError(t, errHttp)

	ts.mockProvider.EXPECT().QueryProviders(gomock.Any()).Return(mockProviderConnectPartialErrorResponse())

	ts.ginRouter.ServeHTTP(ts.recorder, request)

	var actualResponse GetAggregatedSearchLinksResponse
	errUnmarshal := json.Unmarshal(ts.recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, errUnmarshal)

	expHTTPStatus := http.StatusOK
	assert.Equal(t, expHTTPStatus, ts.recorder.Code)

	expResponse := GetAggregatedSearchLinksResponse{
		Data: []Link{
			{
				Url:            "www.yahoo.com/abc6",
				Views:          6000,
				RelevanceScore: 0.6,
			},
			{
				Url:            "www.wikipedia.com/abc1",
				Views:          11000,
				RelevanceScore: 0.1,
			},
		},
		Count: int32(2),
		Errors: []APIError{
			{
				Code:    APIErrorCodeProviderUnavailable,
				Message: "provider google failed",
			},
		},
	}
	assert.Equal(t, expResponse, actualResponse)
}

func Test_GetAggregatedSearchLinks_ProviderConnect_Success_SortByRelevanceScore_LimitResponseHigher(t *testing.T) {
	ts := setupTest(t)
	defer ts.goMockCtrl.Finish()

	request, errHttp := createValidRequest(sortKeyRelevanceScore, "4")
	assert.NoError(t, errHttp)

	ts.mockProvider.EXPECT().QueryProviders(gomock.Any()).Return(mockProviderConnectSuccessResponse())

	ts.ginRouter.ServeHTTP(ts.recorder, request)

	var actualResponse GetAggregatedSearchLinksResponse
	errUnmarshal := json.Unmarshal(ts.recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, errUnmarshal)

	expHTTPStatus := http.StatusOK
	assert.Equal(t, expHTTPStatus, ts.recorder.Code)

	expResponse := GetAggregatedSearchLinksResponse{
		Data: []Link{
			{
				Url:            "www.yahoo.com/abc6",
				Views:          6000,
				RelevanceScore: 0.6,
			},
			{
				Url:            "www.example.com/abc4",
				Views:          4000,
				RelevanceScore: 0.4,
			},
			{
				Url:            "www.wikipedia.com/abc1",
				Views:          11000,
				RelevanceScore: 0.1,
			},
		},
		Count: int32(3),
	}
	assert.Equal(t, expResponse, actualResponse)
}

func Test_GetAggregatedSearchLinks_ProviderConnect_Success_SortByViews_LimitResponseEqual(t *testing.T) {
	ts := setupTest(t)
	defer ts.goMockCtrl.Finish()

	request, errHttp := createValidRequest(sortKeyViews, "2")
	assert.NoError(t, errHttp)

	ts.mockProvider.EXPECT().QueryProviders(gomock.Any()).Return(mockProviderConnectSuccessResponse())

	ts.ginRouter.ServeHTTP(ts.recorder, request)

	var actualResponse GetAggregatedSearchLinksResponse
	errUnmarshal := json.Unmarshal(ts.recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, errUnmarshal)

	expHTTPStatus := http.StatusOK
	assert.Equal(t, expHTTPStatus, ts.recorder.Code)

	expResponse := GetAggregatedSearchLinksResponse{
		Data: []Link{
			{
				Url:            "www.wikipedia.com/abc1",
				Views:          11000,
				RelevanceScore: 0.1,
			},
			{
				Url:            "www.yahoo.com/abc6",
				Views:          6000,
				RelevanceScore: 0.6,
			},
		},
		Count: int32(2),
	}
	assert.Equal(t, expResponse, actualResponse)
}

func Test_GetAggregatedSearchLinks_ProviderConnect_Success_EmptyResponse(t *testing.T) {
	ts := setupTest(t)
	defer ts.goMockCtrl.Finish()

	request, errHttp := createValidRequest(sortKeyViews, "2")
	assert.NoError(t, errHttp)

	ts.mockProvider.EXPECT().QueryProviders(gomock.Any()).Return(mockProviderConnectEmptyResponse())

	ts.ginRouter.ServeHTTP(ts.recorder, request)

	var actualResponse GetAggregatedSearchLinksResponse
	errUnmarshal := json.Unmarshal(ts.recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, errUnmarshal)
	t.Logf("actualResponse: %s", ts.recorder.Body.String())

	expHTTPStatus := http.StatusOK
	assert.Equal(t, expHTTPStatus, ts.recorder.Code)
	assert.Empty(t, actualResponse)
}

func mockProviderConnectEmptyResponse() providerconnect.QueryProvidersResponse {
	return providerconnect.QueryProvidersResponse{
		ProviderResponses: []providerconnect.ProviderResponse{
			{
				Provider: "duckduckgo",
			},
			{
				Provider: "google",
			},
			{
				Provider: "wikipedia",
			},
		},
	}
}

// Mock helper functions

func mockProviderConnectErrorResponse() providerconnect.QueryProvidersResponse {
	return providerconnect.QueryProvidersResponse{
		ProviderResponses: []providerconnect.ProviderResponse{
			{
				Provider: "duckduckgo",
				Error:    fmt.Errorf("duckduckgo failed to retrieve links"),
			},
			{
				Provider: "google",
				Error:    fmt.Errorf("google failed to retrieve links"),
			},
			{
				Provider: "wikipedia",
				Error:    fmt.Errorf("wikipedia failed to retrieve links"),
			},
		},
	}
}

func mockProviderConnectPartialErrorResponse() providerconnect.QueryProvidersResponse {
	return providerconnect.QueryProvidersResponse{
		ProviderResponses: []providerconnect.ProviderResponse{
			mockWikipediaProviderResponse(),
			mockDuckDuckGoProviderResponse(),
			{
				Provider: "google",
				Error:    fmt.Errorf("google failed to retrieve links"),
			},
		},
	}
}

func mockProviderConnectSuccessResponse() providerconnect.QueryProvidersResponse {
	return providerconnect.QueryProvidersResponse{
		ProviderResponses: []providerconnect.ProviderResponse{
			mockWikipediaProviderResponse(),
			mockDuckDuckGoProviderResponse(),
			mockGoogleProviderResponse(),
		},
	}
}

func mockGoogleProviderResponse() providerconnect.ProviderResponse {
	return providerconnect.ProviderResponse{
		Provider: "wikipedia",
		Data: []providerconnect.Link{
			{
				Url:            "www.example.com/abc4",
				Views:          4000,
				RelevanceScore: 0.4,
			},
		},
	}
}

func mockWikipediaProviderResponse() providerconnect.ProviderResponse {
	return providerconnect.ProviderResponse{
		Provider: "wikipedia",
		Data: []providerconnect.Link{
			{
				Url:            "www.wikipedia.com/abc1",
				Views:          11000,
				RelevanceScore: 0.1,
			},
		},
	}
}

func mockDuckDuckGoProviderResponse() providerconnect.ProviderResponse {
	return providerconnect.ProviderResponse{
		Provider: "duckduckgo",
		Data: []providerconnect.Link{
			{
				Url:            "www.yahoo.com/abc6",
				Views:          6000,
				RelevanceScore: 0.6,
			},
		},
	}
}
