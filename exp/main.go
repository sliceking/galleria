package main

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
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

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}

<<<<<<< HEAD
	db.LogMode(true)
	db.AutoMigrate(&User{})
=======
	defer us.Close()
	user, err := us.ByID(2)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
>>>>>>> b77f364dfa117d6c17d39e1785fa63ed6cfbbf32

}
