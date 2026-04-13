package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thedevscott/trug/internal/assert"
)

func TestCommnHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	commonHeaders(next).ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	// Content-Security-Policy check
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, res.Header.Get("Content-Security-Policy"), expectedValue)

	// Referrer-Policy check
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, res.Header.Get("Referrer-Policy"), expectedValue)

	// X-Content-Type-Options check
	expectedValue = "nosniff"
	assert.Equal(t, res.Header.Get("X-Content-Type-Options"), expectedValue)

	// X-Frame-Options check
	expectedValue = "deny"
	assert.Equal(t, res.Header.Get("X-Frame-Options"), expectedValue)

	// X-XSS-Protection check
	expectedValue = "0"
	assert.Equal(t, res.Header.Get("X-XSS-Protection"), expectedValue)

	// Server check
	expectedValue = "Go"
	assert.Equal(t, res.Header.Get("Server"), expectedValue)

	// Next handler called and OK response check
	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
