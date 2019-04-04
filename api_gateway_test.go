package pylon

import (
	"io"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestBuildRequestMethod(test *testing.T) {
	event := &events.APIGatewayProxyRequest{HTTPMethod: "POST"}
	request, err := buildRequestFromGateway(event)
	if err != nil {
		test.Error("Error from buildRequest")
	}
	if request.Method != "POST" {
		test.Error("Invalid Method")
	}
}

func TestBuildRequestUrl(test *testing.T) {
	event := &events.APIGatewayProxyRequest{Path: "/test?key=value&key2=value2"}
	request, err := buildRequestFromGateway(event)
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

func TestBuildRequestHeaders(test *testing.T) {
	headers := make(map[string]string)
	headers["Host"] = "acme.com"
	event := &events.APIGatewayProxyRequest{Path: "/test", Headers: headers}
	request, err := buildRequestFromGateway(event)
	if err != nil {
		test.Error("Error from buildRequest")
	}
	if request.Host != "acme.com" {
		test.Error("Incorrect host header")
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

func TestAPIGatewayResponseWriter(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// A very simple health check.
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		// In the future we could report back on the status of our DB, or our cache
		// (e.g. Redis) by performing a simple PING, and include them in the response.
		io.WriteString(w, `{"alive": true}`)
	})
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := &GatewayResponseWriter{}
	handler.ServeHTTP(res, req)
	if res.response.StatusCode != 200 {
		t.Error("Invalid StatusCode")
	}
	res.finish()

}
