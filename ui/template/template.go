package template

import (
	"github.com/a-h/templ"
	"github.com/timenglesf/personal-site/internal/models"
	"github.com/timenglesf/personal-site/internal/shared"
	"github.com/timenglesf/personal-site/ui/template/components"
	"github.com/timenglesf/personal-site/ui/template/pages"
	"github.com/timenglesf/personal-site/ui/template/partials"
)

type Pages struct {
	Base           func(title string, page templ.Component, data *shared.TemplateData) templ.Component
	PostBase       func(title string, page templ.Component, data *shared.TemplateData) templ.Component
	Index          func(data *shared.TemplateData) templ.Component
	AdminSignup    func(data *shared.TemplateData) templ.Component
	Post           func(data *shared.TemplateData) templ.Component
	CreatePost     func(data *shared.TemplateData) templ.Component
	EditPost       func(data *shared.TemplateData) templ.Component
	AdminLogin     func(data *shared.TemplateData) templ.Component
	AdminDashboard func(data *shared.TemplateData) templ.Component
	About          func(data *shared.TemplateData) templ.Component
}

type Partials struct {
	PageHeader           func(data *shared.TemplateData) templ.Component
	ThemeToggle          func() templ.Component
	AlertError           func(msg, id, tw_classes string) templ.Component
	AlertWarning         func(msg, id, tw_classes string) templ.Component
	AlertSuccess         func(msg, id, tw_classes string) templ.Component
	DashboardBlogPostRow func(p *models.Post) templ.Component
	// Footer      func() templ.Component
}

func CreatePageTemplates() *Pages {
	return &Pages{
		Base:           Base,
		PostBase:       PostBase,
		Index:          pages.Index,
		AdminSignup:    pages.SignUpAdmin,
		AdminLogin:     pages.AdminLogin,
		Post:           pages.Post,
		CreatePost:     pages.CreatePost,
		EditPost:       pages.EditPost,
		AdminDashboard: pages.AdminDashboard,
		About:          pages.AboutPage,
	}
}

func CreatePartialTemplates() *Partials {
	return &Partials{
		PageHeader:           partials.PageHeader,
		ThemeToggle:          partials.ThemeToggle,
		AlertWarning:         components.WarningAlert,
		AlertError:           components.ErrorAlert,
		AlertSuccess:         components.SuccessAlert,
		DashboardBlogPostRow: pages.DashboardBlogPostRow,
		//  Footer:      partials.Footer,
	}
}
