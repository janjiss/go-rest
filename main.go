package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"janjiss.com/rest/users"
)

type CreateUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
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

	r := gin.Default()

	usersService := users.NewUserService(db)

	r.GET("/users", func(c *gin.Context) {
		users := usersService.GetAllUsers()

		c.JSON(http.StatusOK, gin.H{"users": users})
	})

	r.POST("/users", func(c *gin.Context) {
		var userRequest *CreateUser
		var user *users.User
		var err error
		if err = c.ShouldBindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err = usersService.CreateUser(userRequest.Name, userRequest.Email)

		if errors, ok := err.(users.CreateUserError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors.Errors})
			return
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to insert the user into the database", "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	})
	r.Run()
}
