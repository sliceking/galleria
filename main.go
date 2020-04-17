package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sliceking/galleria/controllers"
	"github.com/sliceking/galleria/middleware"
	"github.com/sliceking/galleria/models"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "stanwielga"
	dbname = "galleria_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname,
	)
	services, err := models.NewServices(psqlInfo)
	must(err)

	defer services.Close()
	services.AutoMigrate()

	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	//Image Routes
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	// Gallery Routes
	userMW := middleware.User{UserService: services.User}
	requireUserMW := middleware.RequireUser{User: userMW}
	r.Handle("/galleries",
		requireUserMW.ApplyFn(galleriesC.Index)).Methods("GET")
	r.Handle("/galleries/new",
		requireUserMW.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries",
		requireUserMW.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/edit",
		requireUserMW.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update",
		requireUserMW.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete",
		requireUserMW.ApplyFn(galleriesC.Delete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}",
		galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/images",
		requireUserMW.ApplyFn(galleriesC.Upload)).Methods("POST")
	http.ListenAndServe(":3000", userMW.Apply(r))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
