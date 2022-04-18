package database

import (
	"fmt"
	"log"
	"web/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBconnect *gorm.DB

var err error

func GetDSN(config *config.Config) string {
	// "user=yuanyu password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Taipei"
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort,
	)
}
func DB(config *config.Config) {
	// https://github.com/go-gorm/postgres
	DBconnect, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  GetDSN(config),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}
