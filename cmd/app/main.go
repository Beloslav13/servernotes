package main

import (
	"os"

	"github.com/beloslav13/servernotes/internal/app"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Errorln("No .env file found")
		os.Exit(1)
	}
}

func main() {
	app.Run()
}
