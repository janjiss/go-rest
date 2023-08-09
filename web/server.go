package web

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	auth "janjiss.com/rest/login"
	"janjiss.com/rest/users"
)

func StartServer(db *gorm.DB) {
	us := users.NewUserService(db)

	r := gin.Default()
	r.POST("/login", BuildLoginHandler(us))

	authorized := r.Group("/")
	authorized.Use(auth.JWTAuthMiddleware())
	authorized.GET("/users", BuildGetAllUsersHandler(us))
	authorized.POST("/users", BuildCreateUserHandler(us))

	r.Run()
}
