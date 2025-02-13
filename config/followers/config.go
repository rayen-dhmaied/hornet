package followers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Neo4jURI   string
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
	neo4jURI := os.Getenv("NEO4J_URI")
	dbName := os.Getenv("NEO4J_DB")
	serverPort := os.Getenv("SERVER_PORT")

	// Validate required variables
	if neo4jURI == "" || dbName == "" {
		log.Fatalf("Environment variables MONGO_URI, MONGO_DB, and must be set.")
	}

	// Use default port if not set
	if serverPort == "" {
		serverPort = "8080"
	}

	return &Config{
		Neo4jURI:   neo4jURI,
		DBName:     dbName,
		ServerPort: serverPort,
	}
}
