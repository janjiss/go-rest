package main

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
	"janjiss.com/rest/db"
	"janjiss.com/rest/web"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)

	err := godotenv.Load()
	db := db.NewDB()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	go web.StartServer(db)

	wg.Wait()
}
