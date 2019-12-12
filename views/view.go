package views

import "html/template"

// NewView accepts a list of strings and returns some views, it should
// only be used during setup and not runtime
func NewView(layout string, files ...string) *View {
	files = append(files,
		"views/layout/bootstrap.gohtml",
		"views/layout/navbar.gohtml",
		"views/layout/footer.gohtml",
	)

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
