package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sliceking/galleria/context"
	"github.com/sliceking/galleria/models"
	"github.com/sliceking/galleria/views"
)

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

type Galleries struct {
	New       *views.View
	LoginView *views.View
	gs        models.GalleryService
}

type GalleryForm struct {
	Title string `schema:"title"`
}

//Create is used to make a new gallery
// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	user := context.User(r.Context())
	fmt.Println("create got the user: ", user.ID)
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}

	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	fmt.Fprintln(w, gallery)
}
