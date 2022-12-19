package providerconnect

type Link struct {
	Url            string  `json:"url,omitempty"`
	Views          int32   `json:"views,omitempty"`
	RelevanceScore float32 `json:"relevanceScore,omitempty"`
}

type ProviderResponse struct {
	Provider string
	Data     []Link `json:"data"`
	Error    error
}

type QueryProvidersResponse struct {
	ProviderResponses []ProviderResponse
}
