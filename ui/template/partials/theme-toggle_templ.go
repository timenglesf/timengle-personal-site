// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.680
package partials

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func ThemeToggle() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"p-4\"><label class=\"swap swap-rotate\"><!-- this hidden checkbox controls the state --><input type=\"checkbox\" id=\"theme-toggle\"><!-- sun icon --><svg class=\"swap-on fill-current w-10 h-10\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M5.64 17.36a9 9 0 0012.72 0 9 9 0 000-12.72 9 9 0 00-12.72 0 9 9 0 000 12.72zM12 4a1 1 0 011 1v2a1 1 0 01-2 0V5a1 1 0 011-1zm0 10a1 1 0 011 1v2a1 1 0 01-2 0v-2a1 1 0 011-1zm-7-7a1 1 0 010-1 1 1 0 011-1h2a1 1 0 010 2H6a1 1 0 01-1-1zm10.24 2.24a1 1 0 010-1.42 1 1 0 011.42 0l1.42 1.42a1 1 0 010 1.42 1 1 0 01-1.42 0l-1.42-1.42zM4 12a1 1 0 010-1h2a1 1 0 010 2H5a1 1 0 01-1-1zm13.66 1.34a1 1 0 010 1.42 1 1 0 01-1.42 0l-1.42-1.42a1 1 0 010-1.42 1 1 0 011.42 0l1.42 1.42zM18 6a1 1 0 01-1-1 1 1 0 01-1-1h2a1 1 0 011 1 1 1 0 010 1zM7 19.66a1 1 0 011.42-1.42l1.42 1.42a1 1 0 01-1.42 1.42L7 19.66zm10-2.34h2a1 1 0 010 2h-2a1 1 0 010-2zm-3.66 1.66a1 1 0 010 1.42 1 1 0 01-1.42 0l-1.42-1.42a1 1 0 010-1.42 1 1 0 011.42 0l1.42 1.42zm-6.58-9.24a1 1 0 011.42 0l1.42 1.42a1 1 0 010 1.42 1 1 0 01-1.42 0L6.76 8.58a1 1 0 010-1.42z\"></path></svg><!-- moon icon --><svg class=\"swap-off fill-current w-10 h-10\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path d=\"M12 3.5a9 9 0 000 17.92 9 9 0 0012-12.5 9 9 0 00-12-5.42zm0 16a7 7 0 01-5.32-11.54 7 7 0 009.19 9.19A7 7 0 0112 19.5z\"></path></svg></label><script>\n      const toggleCheckbox = document.getElementById('theme-toggle');\n      toggleCheckbox.addEventListener('change', () => {\n        document.documentElement.classList.toggle('dark', toggleCheckbox.checked);\n      });\n    </script><p class=\"mt-4\">This is some text content.</p></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
