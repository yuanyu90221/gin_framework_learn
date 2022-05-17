package src

import (
	"web/service"

	session "web/middlewares"

	"web/pojo"

	"github.com/gin-gonic/gin"
)

func AddUserRouter(r *gin.RouterGroup) {
	user := r.Group("/users", session.SetSession())
	// user.GET("/", service.FindAllUsers)
	user.GET("/", service.CacheUserAllDecorator(service.RedisAllUser, "user_all", []pojo.User{}))
	// user.GET("/:id", service.FindUserWithId)
	user.GET("/:id", service.CacheOneUseDecorator(service.RedisOneUser, "id", "user_%s", pojo.User{}))
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
