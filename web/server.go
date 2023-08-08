package web

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"janjiss.com/rest/users"
)

func StartServer(db *gorm.DB) {
	us := users.NewUserService(db)

	r := gin.Default()

	r.POST("/login", BuildLoginHandler(us))
	r.GET("/users", BuildGetAllUsersHandler(us))
	r.POST("/users", BuildCreateUserHandler(us))

	r.Run()
}
