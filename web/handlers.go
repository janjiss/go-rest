package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"janjiss.com/rest/helpers/paginator"
	"janjiss.com/rest/users"
)

type CreateUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Login struct {
	Email string `json:"email"`
}

type AllUsers struct {
	Cursor string `json:"cursor"`
}

func BuildCreateUserHandler(us *users.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var userRequest *CreateUser
		var user *users.User
		var err error
		if err = c.ShouldBindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err = us.CreateUser(userRequest.Name, userRequest.Email)

		if errors, ok := err.(users.CreateUserError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors.Errors})
			return
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to insert the user into the database", "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func BuildGetAllUsersHandler(us *users.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {

		var allUsersRequest *AllUsers
		var err error

		if err = c.ShouldBindJSON(&allUsersRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(allUsersRequest.Cursor)

		users, err := us.GetAllUsers(allUsersRequest.Cursor)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}

		var cursor string

		if len(users) > 0 {
			cursor = paginator.EncodeCursor(users[len(users)-1])
		}

		c.JSON(http.StatusOK, gin.H{"users": users, "cursor": cursor})
	}
}

func BuildLoginHandler(us *users.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {

		var loginRequest *Login
		var err error
		var token string

		if err = c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err = us.Login(loginRequest.Email)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
