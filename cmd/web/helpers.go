package main

import (
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/a-h/templ"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
	"github.com/timenglesf/personal-site/internal/models"
	"github.com/timenglesf/personal-site/internal/shared"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) newTemplateData(r *http.Request) shared.TemplateData {
	return shared.TemplateData{
		IsAuthenticated:      app.isAuthenticated(r),
		IsAdmin:              app.isAdmin(r),
		CurrentYear:          time.Now().Year(),
		CSRFToken:            nosurf.Token(r),
		BaseURL:              app.getBaseURLString(r),
		URLPath:              r.URL.Path,
		Flash:                &shared.FlashMessage{},
		MostRecentPublicPost: app.mostRecentPost,
		BlogPosts:            *app.latestPublicPosts,
	}
}

// a helper function to decode form data into a struct using the automatic form decoder
func (app *application) decodeForm(r *http.Request, dst any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	err := app.formDecoder.Decode(&dst, r.PostForm)

	var inValidDecoderError *form.InvalidDecoderError

	if err != nil {
		if errors.As(err, &inValidDecoderError) {
			panic(err)
		}
		return err
	}

	return nil
}

func (app *application) getBaseURLString(r *http.Request) string {
	return r.Proto + "://" + r.Host
}

// returns true if the session data contains the key "authenticatedUserID"
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}

// returns true if the session data contains the key "isAdminRole"
func (app *application) isAdmin(r *http.Request) bool {
	return app.sessionManager.GetBool(r.Context(), "isAdminRole")
}

func (app *application) renderPage(w http.ResponseWriter, r *http.Request, templateFunc func(data *shared.TemplateData) templ.Component, title string, data *shared.TemplateData) {
	page := templateFunc(data)
	base := app.pageTemplates.Base(title, page, data)
	err := base.Render(r.Context(), w)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) renderBlogPostPage(w http.ResponseWriter, r *http.Request, title string, data *shared.TemplateData) {
	page := app.pageTemplates.Post(data)
	base := app.pageTemplates.PostBase("Tim Engle Blog - "+title, page, data)
	err := base.Render(r.Context(), w)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) UpdatePostsOnAppStruct() error {
	mostRecentPublicPost, err := app.post.MostRecentPost(false)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			mostRecentPublicPost = &models.Post{
				Title:   "No posts yet",
				Content: "This is a dummy post that will be removed when the first public post is added.",
			}
		} else {
			return err
		}
	}

	latestPublicPosts, err := app.post.LatestPosts(false)
	if err != nil {
		return err
	}

	if len(latestPublicPosts) == 0 {
		latestPublicPosts = append(latestPublicPosts, *mostRecentPublicPost)
	}

	app.latestPublicPosts = &latestPublicPosts
	app.mostRecentPost = mostRecentPublicPost
	return nil
}
