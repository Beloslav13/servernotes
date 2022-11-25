package database

import (
	"context"
	"database/sql"

	"github.com/beloslav13/servernotes/internal/models"
	"github.com/beloslav13/servernotes/pkg/logger"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Storage struct {
	log logger.Logger
	db  *sql.DB
}

// TODO: Возможно стоит сделать не экспортируемой. и вызывать каждый раз при новом запросе, что бы в конце запросе закрывать соединение
func New(log logger.Logger) (*Storage, error) {
	// connStr := "postgresql://admin:devpass@localhost:5436/servernotes_db?sslmode=disable" // postgresql://localhost:5432/servernotes_db
	connStr := "user=admin password=devpass dbname=servernotes_db sslmode=disable host=db port=5432" // postgresql://localhost:5432/servernotes_db

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Errorln("@@@@@@@@@@", err)
		return nil, err
	}

	// defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{log: log, db: db}, nil
}

func (s *Storage) SaveNotes(ctx context.Context, n *models.Note) error {
	q := `INSERT INTO notes (person_id, category_id, name) VALUES ($1, $2, $3)`

	if _, err := s.db.ExecContext(ctx, q, n.PersonId, n.CategoryId, n.Name); err != nil {
		s.log.Errorf("cannot save note: %w", err)
		return err
	}
	return nil
}
