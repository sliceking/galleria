package views

import "html/template"

// NewView accepts a list of strings and returns some views, it should
// only be used during setup and not runtime
func NewView(files ...string) *View {
	files = append(files, "views/layout/footer.gohtml")

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
	}
}

type View struct {
	Template *template.Template
}
