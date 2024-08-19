package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/timenglesf/personal-site/internal/fileserver"
	"github.com/timenglesf/personal-site/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// fileServer := http.FileServerFS(ui.Files)
	// mux.Handle("/static/", fileServer)

	var fileSvr http.Handler
	if app.cfg.objectStorage.serveStaticObjectStorage {
		fileSvr = &fileserver.ObjectStorageFileServer{ObjectStorageURL: app.cfg.objectStorage.objectStorageURL}
	} else {
		fileSvr = &fileserver.EmbeddedFileServer{FS: http.FS(ui.Files)}
	}

	mux.Handle("/static/", fileSvr)

	dynamic := alice.New(app.sessionManager.LoadAndSave, app.noSurf)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))

	// Posts
	mux.Handle("GET /posts/view/{slug}", dynamic.ThenFunc(app.handleGetBlogPost))
	mux.Handle("GET /posts/latest", dynamic.ThenFunc(app.handleGetLatestBlogPosts))

	adminProtected := dynamic.Append(app.requireAdmin)
	// Protected Post routes

	mux.Handle("POST /posts/create", adminProtected.ThenFunc(app.handleCreateBlogPost))
	mux.Handle("GET /posts/create", adminProtected.ThenFunc(app.handleDisplayCreatePostForm))
	mux.Handle("POST /posts/update/{slug}", adminProtected.ThenFunc(app.handleBlogPostUpdate))
	mux.Handle("GET /posts/edit/{slug}", adminProtected.ThenFunc(app.handleDisplayEditPostForm))
	mux.Handle("POST /posts/edit", adminProtected.ThenFunc(app.handleBlogPostEdit))
	// Admin routes
	//

	mux.Handle("GET /admin/{$}", dynamic.ThenFunc(app.handleDisplayAdminPage))
	mux.Handle("GET /admin/signup", dynamic.ThenFunc(app.handleAdminSignupPage))
	mux.Handle("POST /admin/signup", dynamic.ThenFunc(app.handleAdminSignupPost))
	mux.Handle("GET /admin/login", dynamic.ThenFunc(app.handleAdminLoginPage))
	mux.Handle("POST /admin/login", dynamic.ThenFunc(app.handleAdminLoginPost))
	mux.Handle("POST /admin/logout", dynamic.ThenFunc(app.handleAdminLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, app.commonHeaders)

	return standard.Then(mux)
}
