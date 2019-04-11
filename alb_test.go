package pylon

import (
	"io"
	"log"
	"net/http"
	"testing"
)

func TestALBResponseWriter(t *testing.T) {
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
	res := &ALBResponseWriter{}
	handler.ServeHTTP(res, req)
	if res.response.StatusCode != 200 {
		t.Error("Invalid StatusCode")
	}
	res.finish()
	if res.response.StatusDescription != "200 OK" {
		log.Print(res.response.StatusDescription)
		t.Error("Invalid StatusDescription")
	}

}
