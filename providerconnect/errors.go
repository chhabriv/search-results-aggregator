package providerconnect

import "fmt"

var (
	ErrProviderGoogle     = fmt.Errorf("failed to retrieve links from %s", ProviderGoogle)
	ErrProviderWikipedia  = fmt.Errorf("failed to retrieve links from %s", ProviderWikipedia)
	ErrProviderDuckDuckGo = fmt.Errorf("failed to retrieve links from %s", ProviderDuckDuckGo)
)
