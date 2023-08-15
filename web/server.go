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

	r.GET("/playground", BuildGraphqlPlaygroundHandler())

	r.POST("/login", BuildLoginHandler(us))
	r.POST("/users", BuildCreateUserHandler(us))

	authorized := r.Group("/")
	authorized.Use(auth.JWTAuthMiddleware())
	authorized.GET("/users", BuildGetAllUsersHandler(us))
	authorized.POST("/graphql", BuildGraphqlHandler(us))

	r.Run()
}
