package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/timenglesf/personal-site/internal/models"
	"github.com/timenglesf/personal-site/internal/shared"
	"github.com/timenglesf/personal-site/internal/validator"
)

// Frontmatter struct
type matter struct {
	Title       string   `toml:"title"`
	Date        string   `toml:"date"`
	Description string   `toml:"description"`
	HeaderImage string   `toml:"headerImage"`
	Tags        []string `toml:"tags"`
	Private     bool     `toml:"private"`
}

// Saves a newly created blog to db and redirects to view the post
func (app *application) handleCreateBlogPost(w http.ResponseWriter, r *http.Request) {
	// Parse form
	form, err := app.parseBlogPostForm(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Extract front matter from content
	matter, rest, err := app.extractFrontmatter(form.Content)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	content := string(rest)

	// Validate data
	isDataValid := app.validateBlogPostData(&form, matter, content)
	if !isDataValid {
		data := app.newTemplateData(r)
		data.BlogForm = form
		app.renderPage(w, r, app.pageTemplates.CreatePost, "Create Post", &data)
		return
	}

	// Build post struct
	newPost := app.buildNewPost(form, matter, content, r)

	// Insert post
	id, err := app.post.Insert(*newPost)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Reset mostRecentPublicPost & latestPublicPosts app field
	if err := app.UpdatePostsOnAppStruct(); err != nil {
		app.serverError(w, r, err)
	}

	// Redirect to view of newly created post
	app.redirectToPostView(w, r, id)
}

func (app *application) handleBlogPostEdit(w http.ResponseWriter, r *http.Request) {
	form := struct {
		shared.BlogPostFormData
		ID uint `form:"id"`
	}{}

	if err := app.decodeForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	matter, rest, err := app.extractFrontmatter(form.Content)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	content := string(rest)

	blogForm := shared.BlogPostFormData{
		Title:   matter.Title,
		Content: form.Content,
	}
	isDataValid := app.validateBlogPostData(&blogForm, matter, content)

	if !isDataValid {
		app.rerenderEditPostPage(w, r, blogForm, form.ID)
		return
	}

	// Update the post in the database
	existingPost, err := app.post.GetPostByID(form.ID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	updatedPost := app.buildNewPost(form.BlogPostFormData, matter, content, r)
	updatedPost.ID = form.ID
	updatedPost.CreatedAt = existingPost.CreatedAt

	if err := app.post.Update(updatedPost); err != nil {
		app.serverError(w, r, err)
		return
	}
	// Reset mostRecentPublicPost & latestPublicPosts app field
	if err := app.UpdatePostsOnAppStruct(); err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect to view of updated post
	w.Header().Set("HX-Redirect", fmt.Sprintf("/posts/view/%s", url.QueryEscape(updatedPost.Title)))
	w.WriteHeader(http.StatusSeeOther)
	// http.Redirect(w, r, fmt.Sprintf("/posts/view/%s", url.QueryEscape(post.Title)), http.StatusSeeOther)
}

//////////////////////
/* HELPER FUNCTIONS */
//////////////////////

func (app *application) parseBlogPostForm(r *http.Request) (shared.BlogPostFormData, error) {
	var form shared.BlogPostFormData
	err := app.decodeForm(r, &form)
	return form, err
}

func (app *application) extractFrontmatter(content string) (matter matter, rest []byte, err error,
) {
	rest, err = frontmatter.Parse(strings.NewReader(content), &matter)
	return matter, rest, err
}

func (app *application) validateBlogPostData(form *shared.BlogPostFormData, matter matter, content string) bool {
	form.CheckField(validator.NotBlank(matter.Title), "title", "Title is required")
	form.CheckField(validator.MaxChars(matter.Title, 100), "title", "This field is too long (maximum is 100 characters)")
	form.CheckField(validator.NotBlank(content), "content", "Content is required")

	return form.Valid()
}

func (app *application) buildNewPost(form shared.BlogPostFormData, matter matter, content string, r *http.Request) *models.Post {
	var date time.Time
	var err error
	if matter.Date == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", matter.Date)
		if err != nil {
			date = time.Now()
		}
	}

	userId := app.sessionManager.GetString(r.Context(), sessionUserId)
	return &models.Post{
		Title:       matter.Title,
		Date:        date,
		Description: matter.Description,
		Markdown:    form.Content,
		Content:     content,
		AuthorID:    userId,
		Private:     matter.Private,
		HeaderImage: matter.HeaderImage,
	}
}

// redirectToPostView redirects to the newly created post view.
func (app *application) redirectToPostView(w http.ResponseWriter, r *http.Request, postID uint) {
	createdPost, err := app.post.GetPostByID(postID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flashSuccess", "Post successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/posts/view/%s", url.QueryEscape(createdPost.Title)), http.StatusSeeOther)
}

func (app *application) rerenderEditPostPage(w http.ResponseWriter, r *http.Request, form shared.BlogPostFormData, postID uint) {
	data := app.newTemplateData(r)
	data.BlogForm = form
	data.BlogPost = &models.Post{}
	data.BlogPost.ID = postID
	page := app.pageTemplates.EditPost(&data)
	if err := page.Render(r.Context(), w); err != nil {
		app.serverError(w, r, err)
	}
}
