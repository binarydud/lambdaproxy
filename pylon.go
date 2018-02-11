package pylon

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func buildRequest(ctx *context.Context, event *events.APIGatewayProxyRequest) (*http.Request, error) {
	u, err := url.Parse(event.Path)
	if err != nil {
		return nil, fmt.Errorf("Parse request path: %s", err)
	}
	qs := u.Query()
	for k, v := range event.QueryStringParameters {
		qs.Set(k, v)
	}
	u.RawQuery = qs.Encode()
	rawBody := event.Body
	if event.IsBase64Encoded {
		data, err2 := base64.StdEncoding.DecodeString(rawBody)
		if err2 != nil {
			return nil, fmt.Errorf("Decode base64 request body: %s", err2)
		}
		rawBody = string(data)
	}
	req, err := http.NewRequest(event.HTTPMethod, u.String(), strings.NewReader(rawBody))
	if err != nil {
		return nil, fmt.Errorf("Create request: %s", err)
	}
	for k, v := range event.Headers {
		req.Header.Set(k, v)
	}
	hbody, err := json.Marshal(event.RequestContext)
	if err != nil {
		return nil, fmt.Errorf("Marshal request context: %s", err)
	}
	req.Header.Set("X-ApiGatewayProxy-Context", string(hbody))

	req.Host = event.Headers["Host"]
	return req, nil
}
func ProxyEvent(handler http.Handler) func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		req, err := buildRequest(&ctx, &event)
		if err != nil {
			return events.APIGatewayProxyResponse{}, fmt.Errorf("Build request: %s", err)
		}
		res := &ResponseWriter{}
		handler.ServeHTTP(res, req)
		res.finish()

		return res.response, nil
	}
}
