// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.680
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/timenglesf/personal-site/internal/shared"
	"github.com/timenglesf/personal-site/ui/template/components"
)

func SignUpAdmin(data *shared.TemplateData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"mx-auto w-1/2 max-w-2xl\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if data.Flash.Type != "" {
			templ_7745c5c3_Err = DisplayFlashAlert(data).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form action=\"/admin/signup\" method=\"post\" class=\"\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if data.SignUpForm.FieldErrors["email"] != "" {
			templ_7745c5c3_Err = components.WarningAlert(data.SignUpForm.FieldErrors["email"], "email_warning", "mb-6").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if data.SignUpForm.FieldErrors["confirm_email"] != "" {
			templ_7745c5c3_Err = components.WarningAlert(data.SignUpForm.FieldErrors["confirm_email"], "confirm_email_warning", "mb-6").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if data.SignUpForm.FieldErrors["display_name"] != "" {
			templ_7745c5c3_Err = components.WarningAlert(data.SignUpForm.FieldErrors["display_name"], "display_name_warning", "mb-6").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if data.SignUpForm.FieldErrors["password"] != "" {
			templ_7745c5c3_Err = components.WarningAlert(data.SignUpForm.FieldErrors["password"], "password_warning", "mb-6").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if data.SignUpForm.FieldErrors["confirm_password"] != "" {
			templ_7745c5c3_Err = components.WarningAlert(data.SignUpForm.FieldErrors["confirm_password"], "confirm_password_warning", "mb-6").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex flex-col gap-2 mb-6\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.EmailInput(data.SignUpForm.Email, "email", "Email").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.EmailInput(data.SignUpForm.ConfirmEmail, "confirm_email", "Confirm Email").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.DisplayNameInput(data.SignUpForm.DisplayName, "John Smith").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.PasswordInput("password", "Password").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.PasswordInput("confirm_password", "Confirm Password").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = CSRFHiddenInput(data).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><button class=\"btn btn-primary\" type=\"submit\">Submit</button></form></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
