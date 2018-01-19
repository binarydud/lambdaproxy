package pylon

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

//cribbed from https://github.com/apex/go-apex/blob/415680d65fd80caf4e7da2b79594c11e96347a38/proxy/responsewriter.go
var DefaultTextContentTypes = []string{
	`text/.*`,
	`application/json`,
	`application/.*\+json`,
	`application/xml`,
	`application/.*\+xml`,
}

var textContentTypes []string
var textContentTypesRegexp *regexp.Regexp

func init() {
	err := SetTextContentTypes(DefaultTextContentTypes)
	if err != nil {
		log.Fatal(err)
	}
}

func SetTextContentTypes(types []string) error {
	pattern := "(" + types[0]
	for _, t := range types {
		pattern += "|" + t
	}
	pattern += `)\b.*`

	r, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	textContentTypesRegexp = r
	return nil
}

type ResponseWriter struct {
	response       events.APIGatewayProxyResponse
	output         bytes.Buffer
	headers        http.Header
	headersWritten bool
}

func (w *ResponseWriter) Header() http.Header {
	if w.headers == nil {
		w.headers = make(http.Header)
	}
	return w.headers
}

func (w *ResponseWriter) Write(bs []byte) (int, error) {
	if !w.headersWritten {
		w.WriteHeader(http.StatusOK)
	}
	return w.output.Write(bs)
}

func (w *ResponseWriter) WriteHeader(status int) {
	if w.headersWritten {
		return
	}

	w.response.StatusCode = status

	finalHeaders := make(map[string]string)
	for k, v := range w.headers {
		if len(v) > 0 {
			finalHeaders[k] = v[len(v)-1]
		}
	}

	if value, ok := finalHeaders["Content-Type"]; !ok || value == "" {
		finalHeaders["Content-Type"] = "text/plain; charset=utf-8"
	}

	w.response.Headers = finalHeaders

	w.headersWritten = true
}

// finish writes the accumulated output to the response.Body
func (w *ResponseWriter) finish() {

	// Determine if we should Base64 encode the output
	contentType := w.response.Headers["Content-Type"]

	// Only encode text content types without base64 encoding
	w.response.IsBase64Encoded = !textContentTypesRegexp.MatchString(contentType)

	if w.response.IsBase64Encoded {
		w.response.Body = base64.StdEncoding.EncodeToString(w.output.Bytes())
	} else {
		w.response.Body = w.output.String()
	}
}
