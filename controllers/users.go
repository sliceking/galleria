package controllers

import (
	"net/http"

	"github.com/sliceking/galleria/views"
)

// NewUsers is used to create a new users controller
// this will panic if a template cannot be parsed properly
// and should be used only during initial setup
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

type Users struct {
	NewView *views.View
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}
