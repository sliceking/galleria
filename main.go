package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sliceking/galleria/controllers"
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
	// TODO: Fix this
	// defer us.Close()
	// us.AutoMigrate()

	usersC := controllers.NewUsers(services.User)
	staticC := controllers.NewStatic()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
