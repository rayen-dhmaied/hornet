package followers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Neo4jURI   string
	DBName     string
	User       string
	Password   string
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
	serverPort := os.Getenv("FOLLOWERS_PORT")
	user := os.Getenv("NEO4J_USER")
	password := os.Getenv("NEO4J_PASSWORD")

	// Validate required variables
	if neo4jURI == "" || dbName == "" || user == "" || password == "" {
		log.Fatalf("Environment variables NEO4J_URI, NEO4J_DB, NEO4J_USER and NEO4J_PASSWORD must be set.")
	}

	// Use default port if not set
	if serverPort == "" {
		serverPort = "8080"
	}

	return &Config{
		Neo4jURI:   neo4jURI,
		DBName:     dbName,
		ServerPort: serverPort,
		User:       user,
		Password:   password,
	}
}
