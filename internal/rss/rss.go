package rss

import (
	"encoding/xml"
	"fmt"
	"net/url"

	"github.com/timenglesf/personal-site/internal/models"
)

const (
	//	DOMAIN     = "https://timengle.dev"
	DOMAIN     = "http://localhost:4000"
	POSTS_PATH = "posts/view"
	FEED_PATH  = "rss.xml"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel *Channel `xml:"channel"`
	Atom    string   `xml:"xmlns:atom,attr"`
}

type Channel struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	AtomLink    AtomLink `xml:"atom:link"`
	Description string   `xml:"description"`
	Items       []*Item  `xml:"item"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

func ConvertPostsToRSSItems(posts []models.Post) []*Item {
	items := make([]*Item, 0, len(posts))

	for _, post := range posts {
		if !post.Private {
			u := fmt.Sprintf("%s/%s/%s", DOMAIN, POSTS_PATH, url.QueryEscape(post.Title))
			item := &Item{
				Title:       post.Title,
				Link:        u,
				Description: post.Description,
				PubDate:     post.CreatedAt.Format("Mon, 02 Jan 2006 15:04:05 MST"),
				GUID:        u,
			}
			items = append(items, item)
		}
	}

	return items
}

func CreateRSSStruct(items []*Item) *RSS {
	rss := &RSS{
		Version: "2.0",
		Atom:    "http://www.w3.org/2005/Atom",
		Channel: &Channel{
			Title:       "Tim Engle's Blog",
			Link:        DOMAIN,
			Description: "A blog about software development, technology, and other things.",
			AtomLink: AtomLink{
				Href: fmt.Sprintf("%s/%s", DOMAIN, FEED_PATH),
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Items: items,
		},
	}
	return rss
}

func (rss *RSS) ConvertRSSStructToXML() ([]byte, error) {
	xmlContent, err := xml.MarshalIndent(rss, "", "    ")
	if err != nil {
		return nil, err
	}
	rssFeed := `<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(xmlContent)
	return []byte(rssFeed), nil
}
