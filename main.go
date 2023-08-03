package main

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"janjiss.com/rest/users"
	"janjiss.com/rest/web"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=postgres dbname=gorm port=5432 sslmode=disable",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(users.User{})

	if err != nil {
		fmt.Printf("Err(s):\n%+v\n", err)

		return
	}

	go web.StartServer(db)

	wg.Wait()
}
