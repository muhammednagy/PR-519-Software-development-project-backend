package api

import (
	"fmt"
	"os"
	log "github.com/sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/controllers"
	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/seed"
)

var server = controllers.Server{}
var (
	buildTime string
	version   string
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
	log.Info("Build:", version, buildTime)

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	seed.Load(server.DB)

	server.Run(":8080")

}
