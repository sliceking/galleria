package controllers

import (
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
