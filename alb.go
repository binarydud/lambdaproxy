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

func buildRequestFromALB(event *events.ALBTargetGroupRequest) (*http.Request, error) {
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

// ProxyEvent ...
func ALBProxyEvent(handler http.Handler) func(ctx context.Context, event events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	return func(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
		req, err := buildRequestFromALB(&request)
		if err != nil {
			return events.ALBTargetGroupResponse{}, fmt.Errorf("Build request: %s", err)
		}
		res := &ALBResponseWriter{}
		handler.ServeHTTP(res, req)
		res.finish()

		return res.response, nil
	}
}
