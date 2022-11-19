package transport

import (
	"net/http"
	"time"

	"github.com/beloslav13/servernotes/internal/transport/handlers"
	"github.com/beloslav13/servernotes/pkg/logger"
	"github.com/gorilla/mux"
)

func StartServer(log logger.Logger) {
	router := mux.NewRouter()
	handler := handlers.NewHandler(log)
	handler.Register(router)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln("error start server: %v\n", err)
	}
}
