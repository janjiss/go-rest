package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"janjiss.com/rest/users"
	"janjiss.com/rest/web"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("user=%s dbname=%s port=%s sslmode=%s", os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_SSLMODE")),
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
