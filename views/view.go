package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"
	"github.com/sliceking/galleria/context"
)

var (
	LayoutDir   string = "views/layout/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

// NewView accepts a list of strings and returns some views, it should
// only be used during setup and not runtime
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", errors.New("csrf field is not implemented yet")
		},
	}).ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

// Render is used to render a view with a predefined layout
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
		// do nothing
	default:
		vd = Data{
			Yield: data,
		}
	}

	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})
	if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong. If the problem persists please contact stan.", http.StatusInternalServerError)
		return nil
	}
	io.Copy(w, &buf)
	return nil
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}

	return files
}

func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
