package views

import (
	"html/template"
	"path/filepath"
)

var (
	LayoutDir   string = "views/layout/"
	TemplateExt string = ".gohtml"
)

// NewView accepts a list of strings and returns some views, it should
// only be used during setup and not runtime
func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
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

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}

	return files
}
