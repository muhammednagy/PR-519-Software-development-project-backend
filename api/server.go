package api

import (
	"github.com/joho/godotenv"
	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/controllers"
	log "github.com/sirupsen/logrus"
	"os"
)

var server = controllers.Server{}
var (
	buildTime string
	version   string
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print(".env file found")
	}
}

func Run() {
	var err error = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}
	log.Info("Build:", version, buildTime)
	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	server.Run(":4000")

}
