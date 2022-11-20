package app

import (
	"github.com/beloslav13/servernotes/internal/transport"
	"github.com/beloslav13/servernotes/pkg/logger"
)

func Run() {
	log := logger.GetLogger()
	// s := os.Getenv("TESTIK")
	transport.StartServer(log)
}
