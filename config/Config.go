package config

type Config struct {
	DBUser     string `json:"DBUser"`
	DBPassword string `json:"DBPassword"`
	DBPort     string `json:"DBPort"`
	DBName     string `json:"DBName"`
	DBHost     string `json:"DBHost"`
	Port       string `json:"Port"`
}
