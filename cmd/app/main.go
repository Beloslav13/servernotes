package main

import (
	"github.com/beloslav13/servernotes/internal/app"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	app.Run()
}
