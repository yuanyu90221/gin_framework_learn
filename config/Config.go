package config

import "os"

type Config struct {
	DBUser        string `json:"DBUser"`
	DBPassword    string `json:"DBPassword"`
	DBPort        string `json:"DBPort"`
	DBName        string `json:"DBName"`
	DBHost        string `json:"DBHost"`
	Port          string `json:"Port"`
	RedisPassword string `json:"REDIS_PASSWORD"`
	RedisHost     string `json:"REDIS_HOST"`
	RedisPort     string `json:"REDIS_PORT"`
}

var config = Config{}

func LoadConfig() *Config {
	config.Port = os.Getenv("PORT")
	config.DBPort = os.Getenv("DB_PORT")
	config.DBPassword = os.Getenv("DB_PASSWORD")
	config.DBUser = os.Getenv("DB_USER")
	config.DBName = os.Getenv("DB_NAME")
	config.DBHost = os.Getenv("DB_HOST")
	config.RedisPassword = os.Getenv("REDIS_PASSWORD")
	config.RedisHost = os.Getenv("REDIS_HOST")
	config.RedisPort = os.Getenv("REDIS_PORT")
	return &config
}
func GetConfig() *Config {
	return &config
}
