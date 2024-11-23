package main

import (
	"net/http"

	"github.com/timenglesf/personal-site/internal/rss"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.renderPage(w, r, app.pageTemplates.Index, "Tim Engle - Home", &data)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.renderPage(w, r, app.pageTemplates.About, "Tim Engle - About Me", &data)
}

func (app *application) handleRSSFeed(w http.ResponseWriter, r *http.Request) {
	posts, err := app.post.GetAllPosts(false)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	rssItems := rss.ConvertPostsToRSSItems(posts)
	rssStruct := rss.CreateRSSStruct(rssItems)
	rssFeed, err := rssStruct.ConvertRSSStructToXML()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	_, err = w.Write([]byte(rssFeed))
	app.logger.Error("error writing rss feed writer", "error", err)
}
