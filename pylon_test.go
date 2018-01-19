package pylon

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestBuildRequestMethod(test *testing.T) {
	event := &events.APIGatewayProxyRequest{HTTPMethod: "POST"}
	request, err := buildRequest(event)
	if err != nil {
		test.Error("Error from buildRequest")
	}
	if request.Method != "POST" {
		test.Error("Invalid Method")
	}
}

func TestBuildRequestUrl(test *testing.T) {
	event := &events.APIGatewayProxyRequest{Path: "/test?key=value&key2=value2"}
	request, err := buildRequest(event)
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
	request, err := buildRequest(event)
	if err != nil {
		test.Error("Error from buildRequest")
	}
	if request.Host != "acme.com" {
		test.Error("Incorrect host header")
	}
}
