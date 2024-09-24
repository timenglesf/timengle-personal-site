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
	"math"
	"strconv"
)

func convertIDToString(id uint) string {
	if id > uint(math.MaxInt32) {
		return "0"
	}
	return strconv.Itoa(int(id))
}

func EditPost(data *shared.TemplateData) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<main class=\"container mx-auto h-[calc(100vh-144px)]\"><div class=\"mx-auto w-2/3 h-full flex flex-col\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if data.BlogForm.FieldErrors["title"] != "" {
			templ_7745c5c3_Err = components.WarningAlert(data.BlogForm.FieldErrors["title"], "title_warning", "mb-6").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if data.BlogForm.FieldErrors["content"] != "" && data.BlogForm.FieldErrors["title"] == "" {
			templ_7745c5c3_Err = components.WarningAlert(data.BlogForm.FieldErrors["content"], "content_warning", "mb-6").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form class=\"flex flex-col h-full\" hx-post=\"/posts/edit\" hx-swap=\"innerHTML\" hx-target=\"main\"><div class=\"flex flex-col flex-grow justify-start mb-6\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = labelDisplay("Content:").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.TextAreaInputDisplay(data.BlogForm.Content).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<input type=\"hidden\" name=\"csrf_token\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(data.CSRFToken)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/template/pages/edit_post.templ`, Line: 32, Col: 66}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <input type=\"hidden\" name=\"id\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(convertIDToString(data.BlogPost.ID))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/template/pages/edit_post.templ`, Line: 33, Col: 79}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div><div class=\"mt-4\"><button class=\"btn btn-primary\" type=\"submit\">Update</button> <button hx-get=\"/posts/content\" hx-target=\".modal-box\" hx-swap=\"innerHTML\" hx-include=\"textarea[name=&#39;content&#39;]\" class=\"btn btn-accent text-primary-content\" onclick=\"my_modal_1.showModal()\" type=\"button\">Preview</button></div></form></div></main><dialog id=\"my_modal_1\" class=\"modal\"><div class=\"modal-box font-poppins w-full max-w-5xl\"></div><form method=\"dialog\" class=\"modal-backdrop\"><button>close</button></form></dialog>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
