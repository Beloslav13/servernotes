package postgresql

import (
	"database/sql"
	"github.com/beloslav13/servernotes/pkg/logger"
)

func newConn(log logger.Logger) (*sql.DB, error) {
	// connStr := "postgresql://admin:devpass@localhost:5436/servernotes_db?sslmode=disable" // postgresql://localhost:5432/servernotes_db
	connStr := "user=admin password=devpass dbname=servernotes_db sslmode=disable host=db port=5432"

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Errorln("@@@@@@@@@@", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectDb(log logger.Logger) (*sql.DB, error) {
	db, err := newConn(log)
	if err != nil {
		log.Errorf("cannot connect postgres: %w", err)
		return nil, err
	}
	return db, nil
}
