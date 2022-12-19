package providerconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

var shouldRetryHTTPStatusCodesMap = map[int]bool{
	http.StatusInternalServerError: true,
	http.StatusRequestTimeout:      true,
	http.StatusGatewayTimeout:      true,
	http.StatusServiceUnavailable:  true,
	http.StatusBadGateway:          true,
}

// TODO: backoff considerations for retry
func doHTTPGet(
	ctx context.Context,
	httpClient http.Client,
	url string,
	responseObj interface{},
	retries int,
) error {
	log := log.With().Str("url", url).Logger()
	log.Debug().Msg("doHTTPGet invoked")

	start := time.Now()
	response, httpErr := httpClient.Get(url)
	log.Info().Dur("duration_ms", time.Duration(time.Since(start).Milliseconds())).Msg("provider call duration")
	if httpErr != nil {
		log.Err(httpErr).Msg("http error while calling provider")
		if isHTTPTimeout(httpErr) && retries > 0 {
			log.Info().Msg("http timeout - retrying provider call")
			return doHTTPGet(ctx, httpClient, url, responseObj, retries-1)
		}
		return fmt.Errorf("error calling provider. error: %w", httpErr)
	}
	defer response.Body.Close()
	respBytes, ioError := io.ReadAll(response.Body)
	if ioError != nil {
		return fmt.Errorf("unable to read provider response. error = %w", ioError)
	}
	unmarshalErr := json.Unmarshal(respBytes, responseObj)
	if unmarshalErr != nil {
		log.Err(unmarshalErr).Str("response_body", string(respBytes)).Int("http_status", response.StatusCode).Msg("error unmarshalling provider response")
		return unmarshalErr
	}

	if response.StatusCode != http.StatusOK {
		log.Error().Int("http_status", response.StatusCode).Msg("provider returned non 200 status code")
		if shouldRetryHTTPStatusCodesMap[response.StatusCode] && retries > 0 {
			log.Info().Msg("retrying provider call")
			return doHTTPGet(ctx, httpClient, url, responseObj, retries-1)
		}
		return fmt.Errorf("provider returned non 200 status code. status code = %d", response.StatusCode)
	}

	log.Debug().Msg("doHTTPGet completed")
	return nil
}

func isHTTPTimeout(responseErr error) bool {
	err, ok := responseErr.(net.Error)
	return ok && err.Timeout()
}
