package main

import (
	"fmt"
	"log"
	"os"
	"web/config"
	"web/database"
	"web/src"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var Config = config.Config{}

func main() {
	router := gin.Default()
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
