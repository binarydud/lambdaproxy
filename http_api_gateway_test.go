package pylon

import (
	"io"
	"net/http"
	"testing"
)

func TestHttpBuildRequestMethod(test *testing.T) {
	//event := &events.APIGatewayProxyRequest{HTTPMethod: "POST"}
	http := HttpAPIHttp{Method: "POST"}
	context := HttpAPIGatewayProxyRequestContext{HTTP: http}
	event := &HttpAPIGatewayProxyRequest{RequestContext: context}
	// request, err := buildRequestFromGateway(event)
	request, err := buildRequestFromHttpGateway(event)
	if err != nil {
		test.Error("Error from buildRequest")
	}
	if request.Method != "POST" {
		test.Error("Invalid Method")
	}
}

func TestHttpBuildRequestUrl(test *testing.T) {
	event := &HttpAPIGatewayProxyRequest{Path: "/test?key=value&key2=value2"}
	request, err := buildRequestFromHttpGateway(event)
	if err != nil {
		test.Error("Error from buildRequest")
	}
	if request.URL.Path != "/test" {
		test.Error("Invalid Path")
	}
	if request.URL.RawQuery != "key=value&key2=value2" {
		test.Error("Invalid query string")
	}
}

func TestHttpBuildRequestHeaders(test *testing.T) {
	headers := make(map[string]string)
	headers["Host"] = "acme.com"
	event := &HttpAPIGatewayProxyRequest{Path: "/test", Headers: headers}
	request, err := buildRequestFromHttpGateway(event)
	if err != nil {
		test.Error("Error from buildRequest")
	}
	if request.Host != "acme.com" {
		test.Error("Incorrect host header")
	}
}

func testHttpHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}
