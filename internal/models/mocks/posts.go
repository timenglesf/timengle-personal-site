package mocks

import (
	"time"

	"github.com/timenglesf/personal-site/internal/models"
	"gorm.io/gorm"
)

type PostModel struct{}

func getMockPost(id uint) models.Post {
	mockPosts := []models.Post{
		{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Date(2024, 8, 26, 12, 0, 0, 0, time.UTC),
			},
			Title:       "First Post",
			Date:        time.Date(2024, 8, 26, 12, 0, 0, 0, time.UTC),
			Description: "Wrap up creating a continuous deployment pipeline by integrating Google Cloud's Artifact Registry and Cloud Run.",
			Markdown:    dummyPublicMarkdownText1,
			Content:     dummyContentText1,
			AuthorID:    "1",
			Private:     false,
			HeaderImage: "https://storage.googleapis.com/timengledev-blog-bucket/static/dist/img/icon_sm.png",
		},

		{
			Model: gorm.Model{
				ID:        2,
				CreatedAt: time.Date(2024, 8, 27, 12, 0, 0, 0, time.UTC),
			},
			Title:       "Second Post",
			Date:        time.Date(2024, 8, 27, 12, 0, 0, 0, time.UTC),
			Description: "Another post",
			Markdown:    dummyPublicMarkdownText2,
			Content:     dummyContentText2,
			AuthorID:    "1",
			Private:     false,
			HeaderImage: "https://storage.googleapis.com/timengledev-blog-bucket/static/dist/img/icon_sm.png",
		},

		{
			Model: gorm.Model{
				ID:        3,
				CreatedAt: time.Date(2024, 8, 26, 12, 0, 0, 0, time.UTC),
			},
			Title:       "Third Post",
			Date:        time.Date(2024, 8, 26, 12, 0, 0, 0, time.UTC),
			Description: "Wrap up creating a continuous deployment pipeline by integrating Google Cloud's Artifact Registry and Cloud Run.",
			Markdown:    dummyPrivateMarkdownText1,
			Content:     dummyContentText1,
			AuthorID:    "1",
			Private:     true,
			HeaderImage: "https://storage.googleapis.com/timengledev-blog-bucket/static/dist/img/icon_sm.png",
		},

		{
			Model: gorm.Model{
				ID:        4,
				CreatedAt: time.Date(2024, 8, 27, 12, 0, 0, 0, time.UTC),
			},
			Title:       "Fourth Post",
			Date:        time.Date(2024, 8, 27, 12, 0, 0, 0, time.UTC),
			Description: "Another post",
			Markdown:    dummyPublicMarkdownText2,
			Content:     dummyContentText2,
			AuthorID:    "1",
			Private:     true,
			HeaderImage: "https://storage.googleapis.com/timengledev-blog-bucket/static/dist/img/icon_sm.png",
		},
	}

	return mockPosts[id-1]
}

func (m *PostModel) Insert(p models.Post) (uint, error) {
	return 1, nil
}

func (m *PostModel) Get(id uint) (*models.Post, error) {
	if id > 0 && id < 5 {
		post := getMockPost(id)
		return &post, nil
	}
	return nil, models.ErrNoRecord
}

func (m *PostModel) GetPostByTitle(title string) (*models.Post, error) {
	if title == "First Post" {
		post := getMockPost(1)
		return &post, nil
	}
	if title == "Second Post" {
		post := getMockPost(2)
		return &post, nil
	}
	return nil, models.ErrNoRecord
}

func (m *PostModel) LatestPosts(includePrivatePosts bool) ([]models.Post, error) {
	if includePrivatePosts {
		return []models.Post{getMockPost(1), getMockPost(2)}, nil
	} else {
		return []models.Post{getMockPost(3), getMockPost(4)}, nil
	}
}

func (m *PostModel) MostRecentPost(includePrivatePosts bool) (*models.Post, error) {
	var post models.Post
	if includePrivatePosts {
		post = getMockPost(2)
		return &post, nil
	} else {
		post = getMockPost(4)
		return &post, nil
	}
}

func (m *PostModel) GetPosts(includePrivatePosts bool, page int, pageSize int) ([]models.Post, error) {
	if includePrivatePosts {
		return []models.Post{getMockPost(1), getMockPost(2)}, nil
	}
	return []models.Post{getMockPost(3), getMockPost(4)}, nil
}

func (m *PostModel) Update(p *models.Post) error {
	return nil
}

func (m *PostModel) Count(includePrivatePosts bool) (int64, error) {
	if includePrivatePosts {
		return 4, nil
	}
	return 2, nil
}

func (m *PostModel) GetPostByID(id uint) (*models.Post, error) {
	if id > 0 && id < 5 {
		post := getMockPost(id)
		return &post, nil
	}
	return nil, models.ErrNoRecord
}
