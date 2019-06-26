package service

// Endpoint is a generic API endpoint
type Endpoint struct {
	Resource       string
	HTTPMethod     string
	Route          string
	Description    string
	Explanation    string
	BodyParameters []Parameter
	Requests       []Request
}

// Parameter is a querystring param
type Parameter struct {
	Name          string
	Deprecated    bool
	Description   string
	Type          string
	ValidValues   []string
	ExampleValues []string
}

// Request is an http request/response
type Request struct {
	HTTPMethod     string
	Path           string
	RequestBody    string
	RequestHeaders string
	//RequestHeadersText string
	RequestQueryParameters []Parameter
	//RequestQueryParametersText string
	RequestContentType string
	ResponseStatusCode int
	ResponseStatusText string
	ResponseBody       string
	ResponseHeaders    string
	//ResponseHeadersText string
	ResponseContentType string
	Curl                string
}

// NewEndpoint creates an initialized endpoint
func NewEndpoint(resource, httpMethod, route string) *Endpoint {
	return &Endpoint{
		Resource:       resource,
		HTTPMethod:     httpMethod,
		Route:          route,
		BodyParameters: []Parameter{},
		Requests:       []Request{},
	}
}
