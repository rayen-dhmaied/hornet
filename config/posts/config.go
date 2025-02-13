package posts

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI   string
	DBName     string
	ServerPort string
}

func LoadConfig() *Config {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Warning: .env file not found. Using system environment variables.")
	}

	// Read variables from the environment
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")
	serverPort := os.Getenv("POSTS_PORT")

	// Validate required variables
	if mongoURI == "" || dbName == "" {
		log.Fatalf("Environment variables MONGO_URI, MONGO_DB, and must be set.")
	}

	// Use default port if not set
	if serverPort == "" {
		serverPort = "8080"
	}

	return &Config{
		MongoURI:   mongoURI,
		DBName:     dbName,
		ServerPort: serverPort,
	}
}
