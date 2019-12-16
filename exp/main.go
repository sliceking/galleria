package main

import (
	"fmt"

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

type Order struct {
	gorm.Model
	UserID      uint
	Amount      int
	Description string
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
	db.AutoMigrate(&User{}, &Order{})
	// db.DropTableIfExists(&User{})

	var u User
	if err := db.First(&u).Error; err != nil {
		panic(err)
	}

	createOrder(db, u, 1001, "Fake Desc1")
	createOrder(db, u, 54545, "Fake Desc2")
	createOrder(db, u, 235543453, "Fake Desc3")
	createOrder(db, u, 8877, "Fake Desc4")
}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	err := db.Create(&Order{
		UserID:      user.ID,
		Amount:      amount,
		Description: desc,
	}).Error
	if err != nil {
		panic(err)
	}
}
