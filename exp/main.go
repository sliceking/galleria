package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "stanwielga"
	dbname = "galleria_dev"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
	Color string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname,
	)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.DB().Ping(); err != nil {
		panic(err)
	}
	db.LogMode(true)
	db.AutoMigrate(&User{})
	// db.DropTableIfExists(&User{})

	name, email, color := getInfo()
	u := &User{
		Name:  name,
		Email: email,
		Color: color,
	}
	db.Create(&u)
	fmt.Println(u)
}

func getInfo() (name, email, color string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What is your name?")
	name, _ = reader.ReadString('\n')
	fmt.Println("What is your email?")
	email, _ = reader.ReadString('\n')
	fmt.Println("What is your color?")
	color, _ = reader.ReadString('\n')
	return name, email, color
}
