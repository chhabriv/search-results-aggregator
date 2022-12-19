package api

type Link struct {
	Url            string  `json:"url,omitempty"`
	Views          int32   `json:"views,omitempty"`
	RelevanceScore float32 `json:"relevanceScore,omitempty"`
}

type GetAggregatedSearchLinksResponse struct {
	Data   []Link     `json:"data,omitempty"`
	Count  int32      `json:"count,omitempty"`
	Errors []APIError `json:"errors,omitempty"`
}

type InternalError struct {
	Errors []APIError `json:"errors,omitempty"`
}

type InvalidRequestError struct {
	Errors []BadRequestError `json:"errors,omitempty"`
}
