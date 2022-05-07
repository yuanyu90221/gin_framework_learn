package src

import (
	"web/service"

	session "web/middlewares"

	"github.com/gin-gonic/gin"
)

func AddUserRouter(r *gin.RouterGroup) {
	user := r.Group("/users", session.SetSession())
	user.GET("/", service.FindAllUsers)
	user.GET("/:id", service.FindUserWithId)
	user.POST("/", service.PostUser)
	user.PUT("/:id", service.PutUser)
	user.POST("/login", service.LoginUser)
	user.GET("/check", service.CheckUserSession)
	user.Use(session.AuthSession())
	{
		// delete user
		user.DELETE("/:id", service.DeleteUser)
		// logout user
		user.GET("/logout", service.LogoutUser)
	}
}
