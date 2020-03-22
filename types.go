package pylon

type HttpAPIGatewayProxyRequest struct {
	RouteKey                        string                            `json:"routeKey"` // The resource path defined in API Gateway
	Path                            string                            `json:"rawPath"`  // The url path for the caller
	Headers                         map[string]string                 `json:"headers"`
	QueryStringParameters           map[string]string                 `json:"queryStringParameters"`
	MultiValueQueryStringParameters map[string][]string               `json:"multiValueQueryStringParameters"`
	PathParameters                  map[string]string                 `json:"pathParameters"`
	StageVariables                  map[string]string                 `json:"stageVariables"`
	RequestContext                  HttpAPIGatewayProxyRequestContext `json:"requestContext"`
	Body                            string                            `json:"body"`
	IsBase64Encoded                 bool                              `json:"isBase64Encoded,omitempty"`
	Version                         string                            `json:"version"`
}

type HttpAPIGatewayProxyRequestContext struct {
	AccountID  string      `json:"accountId"`
	APIID      string      `json:"apiId"` // The API Gateway rest API Id
	DomainName string      `json:"domainName"`
	HTTP       HttpAPIHttp `json:"http"`
	RequestID  string      `json:"requestId"`
	RouteKey   string      `json:"routeKey"`
}

type HttpAPIHttp struct {
	Method    string `json:"method"`
	Path      string `json:"path"`
	Protocol  string `json:"protocol"`
	SourceIP  string `json:"sourceIp"`
	UserAgent string `json:"userAgent"`
}
