package transport

import (
	"net/http"
	"time"

	"github.com/beloslav13/servernotes/internal/database"
	"github.com/beloslav13/servernotes/internal/handlers"
	"github.com/beloslav13/servernotes/pkg/logger"
	"github.com/gorilla/mux"
)

func StartServer(log logger.Logger, db *database.Storage) {
	router := mux.NewRouter()
	handler := handlers.NewHandler(log, db)
	handler.Register(router)
	log.Infoln("Start server")

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln("error start server: %v\n", err)
	}
}
