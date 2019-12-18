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
	defer us.Close()
	us.DestructiveReset()

	u := models.User{
		Name:  "superasdfs man",
		Email: "super@man.com",
	}

	if err := us.Create(&u); err != nil {
		panic(err)
	}

	// u.Email = "paper@person.com"
	// if err = us.Update(&u); err != nil {
	// 	panic(err)
	// }

	if err := us.Delete(u.ID); err != nil {
		panic(err)
	}

	userByID, err := us.ByID(u.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(userByID)
}
