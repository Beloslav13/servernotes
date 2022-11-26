package database

import (
	"context"
	"database/sql"
	"github.com/beloslav13/servernotes/internal/models"
	"github.com/beloslav13/servernotes/pkg/logger"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var log = logger.GetLogger()

func newConn() (*sql.DB, error) {
	// connStr := "postgresql://admin:devpass@localhost:5436/servernotes_db?sslmode=disable" // postgresql://localhost:5432/servernotes_db
	connStr := "user=admin password=devpass dbname=servernotes_db sslmode=disable host=db port=5432" // postgresql://localhost:5432/servernotes_db

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

func SaveNotes(ctx context.Context, n *models.Note) error {
	db, err := newConn()
	defer db.Close()
	if err != nil {
		log.Errorf("cannot connect database: %w", err)
		return err
	}

	q := `INSERT INTO notes (person_id, category_id, name) VALUES ($1, $2, $3)`

	if _, err := db.ExecContext(ctx, q, n.PersonId, n.CategoryId, n.Name); err != nil {
		log.Errorf("cannot save note: %w", err)
		return err
	}
	return nil
}
