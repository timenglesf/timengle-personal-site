package shared

import (
	"context"
	"io"
	"time"

	"github.com/a-h/templ"
	"github.com/timenglesf/personal-site/internal/models"
	"github.com/timenglesf/personal-site/internal/validator"
)

const (
	DateLayout     = "January 2, 2006"
	DateTimeLayout = "January 2, 2006 at 3:04 PM"
)

type FlashMessage struct {
	Message string
	Type    string
}

type TemplateData struct {
	BlogPosts            []models.Post
	CSRFToken            string
	BaseURL              string
	User                 models.User
	BlogForm             BlogPostFormData
	SignUpForm           AdminSignUpForm
	LoginForm            AdminLoginForm
	Flash                *FlashMessage
	BlogPost             *models.Post
	MostRecentPublicPost *models.Post
	CurrentYear          int
	TotalPostCount       int
	CurrentPage          int
	IsAuthenticated      bool
	IsAdmin              bool
}

// Form data

type AdminSignUpForm struct {
	validator.Validator `form:"-"`
	Email               string `form:"email"`
	ConfirmEmail        string `form:"confirm_email"`
	Password            string `form:"password"`
	ConfirmPassword     string `form:"confirm_password"`
	DisplayName         string `form:"display_name"`
}

type AdminLoginForm struct {
	validator.Validator `form:"-"`
	Email               string `form:"email"`
	Password            string `form:"password"`
}

type BlogPostFormData struct {
	validator.Validator `form:"-"`
	Title               string `form:"title"`
	Content             string `form:"content"`
	// Tags
}

// Converts html string to a templ.Component
func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}

func HumanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DateLayout)
}

func HumanDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DateTimeLayout)
}
