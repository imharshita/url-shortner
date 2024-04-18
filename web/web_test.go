package web_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imharshita/url-shortner/short"
	"github.com/imharshita/url-shortner/web/api"
)

func TestWeb(t *testing.T) {
	// short service
	short.Start()

	// Test ShortAndExpandURL
	t.Run("ShortAndExpandURL", func(t *testing.T) {
		// Short URL test
		payload := []byte(`{"longURL":"https://example.com"}`)
		req, err := http.NewRequest("POST", "/short", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(api.ShortURL)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		expected := `{"shortURL":"c984d06"}`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}

		// Expand URL test
		payload = []byte(`{"shortURL":"c984d06"}`)
		req, err = http.NewRequest("POST", "/expand", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()
		handler = http.HandlerFunc(api.ExpandURL)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		expected = `{"longURL":"https://example.com"}`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	})
}
