package template

import (
	"testing"
)

// TestCreatePageTemplates ensures that all fields in the Pages struct are initialized.
func TestCreatePageTemplates(t *testing.T) {
	pages := CreatePageTemplates()

	tests := []struct {
		name string
		fn   interface{}
	}{
		{"Base", pages.Base},
		{"PostBase", pages.PostBase},
		{"Index", pages.Index},
		{"AdminSignup", pages.AdminSignup},
		{"Post", pages.Post},
		{"CreatePost", pages.CreatePost},
		{"EditPost", pages.EditPost},
		{"AdminLogin", pages.AdminLogin},
		{"AdminDashboard", pages.AdminDashboard},
	}

	for _, tt := range tests {
		if tt.fn == nil {
			t.Errorf("Field %s in Pages struct is not initialized", tt.name)
		}
	}
}

// TestCreatePartialTemplates ensures that all fields in the Partials struct are initialized.
func TestCreatePartialTemplates(t *testing.T) {
	partials := CreatePartialTemplates()

	tests := []struct {
		name string
		fn   interface{}
	}{
		{"PageHeader", partials.PageHeader},
		{"ThemeToggle", partials.ThemeToggle},
		{"AlertError", partials.AlertError},
		{"AlertWarning", partials.AlertWarning},
		{"AlertSuccess", partials.AlertSuccess},
		{"DashboardBlogPostRow", partials.DashboardBlogPostRow},
		// Uncomment if Footer is added back
		// {"Footer", partials.Footer},
	}

	for _, tt := range tests {
		if tt.fn == nil {
			t.Errorf("Field %s in Partials struct is not initialized", tt.name)
		}
	}
}
