package models

import (
	"errors"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title      string     `gorm:"not null:unique"`
	Content    string     `gorm:"not null"`
	AuthorID   string     `gorm:"not null"`
	Author     User       `gorm:"foreignKey:AuthorID"`
	Tags       []Tag      `gorm:"many2many:post_tags"`
	Categories []Category `gorm:"many2many:post_categories"`
	Private    bool       `gorm:"not null;default:false"`
}

type PostModel struct {
	DB *gorm.DB
}

func (m *PostModel) Insert(title, content string, private bool, authorID string) (uint, error) {
	p := Post{
		Title:    title,
		Content:  content,
		Private:  private,
		AuthorID: authorID,
	}
	result := m.DB.Create(&p)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return 0, ErrDuplicateTitle
		}
		return 0, result.Error
	}
	return p.ID, nil
}

func (m *PostModel) Get(id uint) (*Post, error) {
	var p Post
	result := m.DB.Preload("Tags").Preload("Categories").First(&p, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNoRecord
		}
		return nil, result.Error
	}
	return &p, nil
}

func (m *PostModel) GetPostByTitle(title string) (*Post, error) {
	var p Post
	result := m.DB.Preload("Tags").Preload("Categories").Where("title = ?", title).First(&p)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNoRecord
		}
		return nil, result.Error
	}
	return &p, nil
}

func (m *PostModel) LatestPosts(includePrivatePosts bool) ([]Post, error) {
	var posts []Post
	result := m.DB.Preload("Tags").Preload("Categories").Where("private = ?", includePrivatePosts).Order("created_at DESC").Limit(10).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func (m *PostModel) MostRecentPost(includePrivatePosts bool) (*Post, error) {
	var post Post
	result := m.DB.Preload("Tags").Preload("Categories").Where("private = ?", includePrivatePosts).Order("created_at DESC").First(&post)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		}
		return nil, result.Error
	}
	return &post, nil
}

func (m *PostModel) GetPosts(includePrivatePosts bool, page int, pageSize int) ([]Post, error) {
	var posts []Post
	offset := (page - 1) * pageSize
	result := m.DB.Preload("Tags").Preload("Categories").Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func (m *PostModel) Count(includePrivatePosts bool) (int64, error) {
	var count int64
	result := m.DB.Model(&Post{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (m *PostModel) Update(p *Post) error {
	result := m.DB.Save(p)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *PostModel) GetPostByID(id uint) (*Post, error) {
	var p Post
	result := m.DB.Preload("Tags").Preload("Categories").First(&p, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNoRecord
		}
		return nil, result.Error
	}
	return &p, nil
}
