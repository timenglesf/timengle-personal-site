package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/justinas/nosurf"
	"github.com/timenglesf/personal-site/internal/models"
	"github.com/timenglesf/personal-site/internal/shared"
	"github.com/timenglesf/personal-site/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) handleDisplayAdminPage(w http.ResponseWriter, r *http.Request) {
	// If there is no admin user in the database, display the admin signup page
	_, err := app.user.GetAdmin()
	if err != nil {
		// Display the admin signup page
		if errors.Is(err, models.ErrNoRecord) {
			data := app.newTemplateData(r)
			app.renderPage(w, r, app.pageTemplates.AdminSignup, "Admin Sign Up", &data)
			return

		}
		// Handle other errors
		// TODO: Display an error message on the page using HTMX
		app.serverError(w, r, err)
	}

	data := app.newTemplateData(r)
	// display admin login page
	if !app.isAdmin(r) {
		app.renderPage(w, r, app.pageTemplates.AdminLogin, "Admin Login", &data)
		return
	}

	recentPosts, err := app.post.GetPosts(true, 1, 10)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	totalPosts, err := app.post.Count(true)
	if err != nil {
		app.serverError(w, r, err)
	}

	data.BlogPosts = recentPosts
	data.TotalPostCount = int(totalPosts)
	data.CurrentPage = 1
	app.renderPage(w, r, app.pageTemplates.AdminDashboard, "Dashboard", &data)
}

func (app *application) handleAdminSignupPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.renderPage(w, r, app.pageTemplates.AdminSignup, "Admin Sign Up", &data)
}

func (app *application) handleAdminSignupPost(w http.ResponseWriter, r *http.Request) {
	// parse and validate form
	var form shared.AdminSignUpForm

	err := app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.Email = strings.TrimSpace(form.Email)
	form.ConfirmEmail = strings.TrimSpace(form.ConfirmEmail)

	form.CheckField(validator.NotBlank(form.Email), "email", "Email is required")
	form.CheckField(validator.ValidEmail(form.Email), "email", "Invalid email format")
	form.CheckField(validator.MaxChars(form.Email, 100), "email", "Email is too long (maximum is 100 characters)")
	form.CheckField(validator.NotBlank(form.ConfirmEmail), "confirm_email", "Confirm Email is required")
	form.CheckField(validator.EqualEmails(form.Email, form.ConfirmEmail), "confirm_email", "Emails do not match")

	form.CheckField(validator.NotBlank(form.DisplayName), "display_name", "Name is required")
	form.CheckField(validator.MaxChars(form.DisplayName, 50), "display_name", "Name is too long (maximum is 50 characters)")
	form.CheckField(validator.MinChars(form.DisplayName, 2), "display_name", "Name is too short (minimum is 2 characters)")

	form.CheckField(validator.NotBlank(form.Password), "password", "Password is required")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be at least 8 characters long")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "Password is too long (maximum is 100 characters)")
	form.CheckField(validator.NotBlank(form.ConfirmPassword), "confirm_password", "Confirm Password is required")
	form.CheckField(validator.EqualStrings(form.Password, form.ConfirmPassword), "confirm_password", "Passwords do not match")

	data := app.newTemplateData(r)
	data.SignUpForm = form

	if !form.Valid() {
		app.renderPage(w, r, app.pageTemplates.AdminSignup, "Admin Signup", &data)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 12)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Get an admin if one exists
	adminData, err := app.user.GetAdmin()
	if err != nil {
		// If the error is not a no rows request than this error is unexpected
		if !errors.Is(err, models.ErrNoRecord) {
			app.serverError(w, r, err)
			return
		}
	}

	// If an admin already exists, redirect
	if adminData != nil {
		app.sessionManager.Put(r.Context(), "flashError", "Creating admin")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	_, err = app.user.Insert(form.DisplayName, form.Email, string(hashedPassword), true)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateAdmin) {
			app.sessionManager.Put(r.Context(), "flashError", "Creating admin")
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		// If there is no record we want to continue and create one
		if !errors.Is(err, models.ErrNoRecord) {
			app.serverError(w, r, err)
			return
		}
	}

	app.sessionManager.Put(r.Context(), "flashSuccess", "Admin account created successfully")
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

func (app *application) handleAdminLoginPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	flashSuccess := app.sessionManager.PopString(r.Context(), "flashSuccess")
	if flashSuccess != "" {
		data.Flash = &shared.FlashMessage{Message: flashSuccess, Type: "success"}
	}

	flashError := app.sessionManager.PopString(r.Context(), "flashError")
	if flashError != "" {
		data.Flash = &shared.FlashMessage{Message: flashError, Type: "error"}
	}

	app.renderPage(w, r, app.pageTemplates.AdminLogin, "Admin Login", &data)
}

func displayAdminLoginWithInvalidCredAlert(app *application, w http.ResponseWriter, r *http.Request, data *shared.TemplateData) {
	data.Flash = &shared.FlashMessage{Message: "Email or Password Incorrect", Type: "warning"}
	w.WriteHeader(http.StatusUnprocessableEntity)
	app.renderPage(w, r, app.pageTemplates.AdminLogin, "Admin Login", data)
}

func (app *application) handleAdminLoginPost(w http.ResponseWriter, r *http.Request) {
	var form shared.AdminLoginForm

	err := app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.Email = strings.TrimSpace(form.Email)

	form.CheckField(validator.NotBlank(form.Email), "email", "Email is required")
	form.CheckField(validator.ValidEmail(form.Email), "email", "Invalid email format")
	form.CheckField(validator.MaxChars(form.Email, 100), "email", "Email is too long (maximum is 100 characters)")

	form.CheckField(validator.NotBlank(form.Password), "password", "Password is required")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be at least 8 characters long")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "Password is too long (maximum is 100 characters)")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.LoginForm = form
		app.renderPage(w, r, app.pageTemplates.AdminLogin, "Admin Login", &data)
		return
	}

	if _, err = app.user.GetAdmin(); err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			data := app.newTemplateData(r)
			data.Flash = &shared.FlashMessage{Message: "Admin does not exist", Type: "warning"}
			w.WriteHeader(http.StatusUnprocessableEntity)
			app.renderPage(w, r, app.pageTemplates.AdminSignup, "Sign Up", &data)
			return
		}
	}

	data := app.newTemplateData(r)
	data.LoginForm = form

	user, err := app.user.Authenticate(form.Email, form.Password, true)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			displayAdminLoginWithInvalidCredAlert(app, w, r, &data)
			return
		}
		if errors.Is(err, models.ErrInvalidCredentials) {
			displayAdminLoginWithInvalidCredAlert(app, w, r, &data)
			return
		}
		if errors.Is(err, models.ErrUpdateUser) {
			app.logger.Error("Error updating user", "error", err)
		}
	}
	if err = app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", user.ID)
	app.sessionManager.Put(r.Context(), "isAdminRole", true)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (app *application) handleAdminLogoutPost(w http.ResponseWriter, r *http.Request) {
	sentCSRFToken, err := r.Cookie("csrf_token")
	if err != nil {
		app.clientError(w, http.StatusForbidden)
		return
	}

	if !nosurf.VerifyToken(nosurf.Token(r), sentCSRFToken.Value) {
		alertComponent := app.partialTemplates.AlertError("CSRF token incorrect. Please try again", "", "error-message container max-w-screen-sm mx-auto mb-6")
		if err = alertComponent.Render(r.Context(), w); err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, r, err)
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Remove(r.Context(), "isAdminRole")

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusNoContent)
}
