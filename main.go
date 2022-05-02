package main

import (
	"fmt"
	"io"
	"log"
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

var Config = config.Config{}

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
	Config.Port = os.Getenv("PORT")
	Config.DBPort = os.Getenv("DB_PORT")
	Config.DBPassword = os.Getenv("DB_PASSWORD")
	Config.DBUser = os.Getenv("DB_USER")
	Config.DBName = os.Getenv("DB_NAME")
	Config.DBHost = os.Getenv("DB_HOST")
	log.Printf("%v", Config)
	go func() {
		database.DB(&Config)
	}()
	v1 := router.Group("/v1")
	src.AddUserRouter(v1)
	router.Run(fmt.Sprintf(":%s", Config.Port))
}
