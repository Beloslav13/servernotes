package transport

import (
	"github.com/beloslav13/servernotes/internal/notes"
	"github.com/beloslav13/servernotes/internal/person"
	"net/http"
	"time"

	"github.com/beloslav13/servernotes/internal/handlers"
	"github.com/beloslav13/servernotes/pkg/logger"
	"github.com/gorilla/mux"
)

func StartServer(log logger.Logger) {
	router := mux.NewRouter()
	noteRepository := notes.NewRepository(log)
	personRepository := person.NewRepository(log)
	noteHandler := handlers.NewNoteHandler(noteRepository, log)
	personHandler := handlers.NewPersonHandler(personRepository, log)
	noteHandler.Register(router)
	log.Infoln("create Note handler...")
	personHandler.Register(router)
	log.Infoln("create Person handler...")
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
