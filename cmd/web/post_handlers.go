package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/justinas/nosurf"
	"github.com/timenglesf/personal-site/internal/models"
	"github.com/timenglesf/personal-site/internal/shared"
	"github.com/yuin/goldmark"

	figure "github.com/mangoumbrella/goldmark-figure"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

// Create Post Handlers

// Displays form to create a blog post
func (app *application) handleDisplayCreatePostForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	if !data.IsAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.renderPage(w, r, app.pageTemplates.CreatePost, "Create Post", &data)
}

// Render blog Post by Title
func (app *application) handleGetBlogPost(w http.ResponseWriter, r *http.Request) {
	titleId := r.PathValue("slug")
	targetPostTitle, err := url.QueryUnescape(titleId)
	if err != nil {
		fmt.Println(err)
		app.clientError(w, http.StatusBadRequest)
	}

	post, err := app.post.GetPostByTitle(targetPostTitle)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Parse front matter from content

	// Reject unauthorized access to private posts
	if post.Private {
		if !app.isAdmin(r) {
			app.logger.Warn("unauthorized access to url", "url", r.URL.Path, "ip", r.RemoteAddr)
			referer := r.Referer()
			if referer != "" {
				http.Redirect(w, r, referer, http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			return
		}
	}

	data := app.newTemplateData(r)
	data.BlogPost = post

	contentHTML, err := convertMarkdownContentToHTML(post.Content)
	if err != nil {
		app.logger.Error("Error converting markdown to html", "error", err)
		app.serverError(w, r, err)
		return
	}

	data.BlogPost.Content = contentHTML

	// Get flash message from session
	flashSuccess := app.sessionManager.PopString(r.Context(), "flashSuccess")
	if flashSuccess != "" {
		data.Flash = &shared.FlashMessage{Message: flashSuccess, Type: "Post created successfully"}
	}

	app.renderBlogPostPage(w, r, post.Title, &data)
	// app.renderPage(w, r, app.pageTemplates.Post, "Post", &data)
}

func (app *application) handleGetLatestBlogPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := app.post.LatestPosts(false)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	for _, post := range posts {
		fmt.Fprintf(w, "%+v\n", post)
	}
}

func (app *application) handleBlogPostUpdate(w http.ResponseWriter, r *http.Request) {
	sentCSRFTOKEN, err := r.Cookie("csrf_token")
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if !nosurf.VerifyToken(nosurf.Token(r), sentCSRFTOKEN.Value) {
		return
	}

	slug := r.PathValue("slug")
	query := r.URL.Query()

	title, _ := url.QueryUnescape(slug)

	post, err := app.post.GetPostByTitle(title)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	for key, value := range query {
		switch key {
		case "title":
			post.Title = value[0]
		case "content":
			post.Content = value[0]
		case "private":
			post.Private = !post.Private
		}
	}

	if err := app.post.Update(post); err != nil {
		app.serverError(w, r, err)
		return
	}

	// Reset mostRecentPublicPost & latestPublicPosts app field
	if err := app.UpdatePostsOnAppStruct(); err != nil {
		app.serverError(w, r, err)
	}

	// Send updated blog post row if this is an updated to the private column
	if query.Get("private") != "" {
		updatedRowComponenet := app.partialTemplates.DashboardBlogPostRow(post)
		if err = updatedRowComponenet.Render(r.Context(), w); err != nil {
			app.serverError(w, r, err)
		}
	}
}

func (app *application) handleDisplayEditPostForm(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	id, err := strconv.Atoi(slug)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if id < 0 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	post, err := app.post.GetPostByID(uint(id))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var form shared.BlogPostFormData
	form.Title = post.Title
	form.Content = post.Markdown

	templateData := app.newTemplateData(r)
	templateData.BlogForm = form
	templateData.BlogPost = post

	page := app.pageTemplates.EditPost(&templateData)
	if err = page.Render(r.Context(), w); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) handleGetBlogModal(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	c := vals.Get("content")

	matter, contentByte, err := app.extractFrontmatter(c)
	if err != nil {
		app.logger.Error("Error extracting front matter", "error", err)
		fmt.Fprintf(w, "There was an error previewing your content")
	}

	header := `<h1 class="font-poppins text-3xl">` + matter.Title + `</h1>`
	contentHTML, err := convertMarkdownContentToHTML(string(contentByte))
	if err != nil {
		app.logger.Error("Error converting markdown to html", "error", err)
		fmt.Fprintf(w, "There was an error previewing your content")
	}

	proseWrapper := `<div class="
			prose 
      prose-sm
      prose-code:font-sans
      prose-pre:font-sans
      prose-headings:font-sans
      hover:prose-a:text-info
      prose-a:duration-300
      lg:prose-lg
">` + contentHTML + `</div>`

	fmt.Fprintf(w, "%s", header+proseWrapper)
}

func convertMarkdownContentToHTML(content string) (string, error) {
	mdRenderer := goldmark.New(
		goldmark.WithExtensions(
			figure.Figure,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
			),
		),
	)
	var buf bytes.Buffer
	if err := mdRenderer.Convert([]byte(content), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
