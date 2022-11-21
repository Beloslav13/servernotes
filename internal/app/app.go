package app

import (
	"github.com/beloslav13/servernotes/internal/database"
	"github.com/beloslav13/servernotes/internal/transport"
	"github.com/beloslav13/servernotes/pkg/logger"
	
)

func Run() {
	log := logger.GetLogger()
	log.Infoln("Start app...")
	db, err := database.New(log)
	if err != nil {
		log.Fatalln(err)
	}
	log.Infoln("connect database success...")

	transport.StartServer(log, db)
}
