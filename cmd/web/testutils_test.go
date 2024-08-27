package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/timenglesf/personal-site/internal/models/mocks"
	"github.com/timenglesf/personal-site/ui/template"
)

func newTestApplication(t *testing.T) *application {
	cfg := &config{
		environment: "test",
	}
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	post := &mocks.PostModel{}
	latestPublicPosts, _ := post.LatestPosts(false)
	return &application{
		// logger:           slog.New(slog.NewTextHandler(io.Discard, nil)),
		logger:            slog.New(slog.NewTextHandler(os.Stdout, nil)),
		cfg:               cfg,
		user:              &mocks.UserModel{},
		post:              &mocks.PostModel{},
		pageTemplates:     template.CreatePageTemplates(),
		partialTemplates:  template.CreatePartialTemplates(),
		formDecoder:       formDecoder,
		sessionManager:    sessionManager,
		latestPublicPosts: &latestPublicPosts,
		mostRecentPost:    &latestPublicPosts[1],
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
