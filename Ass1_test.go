package main

import "testing"
import "net/http"
import "net/http/httptest"

func Test_handler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerFunc))
	defer ts.Close()

	res, err := http.Get(ts.URL + "github.com/apache/kafka")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, res.StatusCode)
		return
	}
}