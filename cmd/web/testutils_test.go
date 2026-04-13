package main

import (
	"bytes"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/thedevscott/trug/internal/models/mocks"
)

// newTestApplication returns an instance of our application struct containing
// mocked dependencies.
func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		logger:         slog.New(slog.DiscardHandler),
		users:          &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

// testServer embeds an httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// newTestServer initializes and returns a new instance of our custom testServer
// type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	// stop test server from followign redirects
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		// tells client to stop and return the received response
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// resetClientCookieJar resets the test server client to use a new and empty
// cookie jar.
func (ts *testServer) resetClientCookieJar(t *testing.T) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	ts.Client().Jar = jar
}

// testResponse a struct to hold data about responses from the test server. Note
// that this struct includes fields for the HTTP response headers and cookies,
// as well as the status code and body.
type testResponse struct {
	status  int
	headers http.Header
	cookies []*http.Cookie
	body    string
}

// get is a method on our custom testServer type. This makes a GET request to a
// given url path using the test server client and it returns a  testResponse
// struct containing the response data.
func (ts *testServer) get(t *testing.T, urlPath string) testResponse {
	req, err := http.NewRequest(http.MethodGet, ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return testResponse{
		status:  res.StatusCode,
		headers: res.Header,
		cookies: res.Cookies(),
		body:    string(bytes.TrimSpace(body)),
	}
}

// extractCSRFToken returns the form CSRF token required to make a valid post
func extractCSRFToken(t *testing.T, body string) string {
	csrfTokenRX := regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(matches[1])
}

// postForm creates a method for sending POST requests to the test server.
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) testResponse {
	req, err := http.NewRequest(http.MethodPost, ts.URL+urlPath, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return testResponse{
		status:  res.StatusCode,
		headers: res.Header,
		cookies: res.Cookies(),
		body:    string(bytes.TrimSpace(body)),
	}
}
