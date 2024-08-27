package main

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/timenglesf/personal-site/internal/assert"
)

func TestHandleGetBlogPost(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "validURLPath",
			urlPath:  fmt.Sprintf("/posts/view/%s", url.QueryEscape("First Post")),
			wantCode: http.StatusOK,
			wantBody: "First Post",
		},
		{
			name:     "secondValidURLPath",
			urlPath:  fmt.Sprintf("/posts/view/%s", url.QueryEscape("Second Post")),
			wantCode: http.StatusOK,
			wantBody: "Second Post",
		},
		{
			name:     "nonExistentPost",
			urlPath:  "/posts/view/does-not-exist",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "digitId",
			urlPath:  "/posts/view/123",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			fmt.Print(tt.urlPath)
			assert.Equal(t, code, tt.wantCode)
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
