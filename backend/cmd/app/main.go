// @title EduApp API
// @version 1.0
// @description REST API for EduApp
// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"flag"
	"log"
	"gitlab.com/w0ikid/study-platform/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	// configFile := flag.String("config", "configs/.env", "Path to configuration file")
	configFile := flag.String("config", "/app/configs/.env", "Path to configuration file")

	flag.Parse()

	err := godotenv.Load(*configFile)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// app
	if err := app.Run(*configFile); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}