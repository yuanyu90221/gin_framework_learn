package main

import (
	"fmt"
	"io"
	"os"
	"web/config"
	"web/database"
	"web/middlewares"
	"web/src"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
)

func setupLogging() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {
	setupLogging()
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("userpasswd", middlewares.UserPassswd)
	}
	router.Use(gin.Recovery(), gin.BasicAuth(
		gin.Accounts{os.Getenv("BASIC_AUTH_USER"): os.Getenv("BASIC_AUTH_PASSWORD")}),
		middlewares.Logger())
	Config := config.LoadConfig()
	go func() {
		database.DB(Config)
	}()
	v1 := router.Group("/v1")
	src.AddUserRouter(v1)
	router.Run(fmt.Sprintf(":%s", Config.Port))
}
