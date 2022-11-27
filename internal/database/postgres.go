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

func connectDb() (*sql.DB, error) {
	db, err := newConn()
	if err != nil {
		log.Errorf("cannot connect database: %w", err)
		return nil, err
	}
	return db, nil
}

func CreateNote(ctx context.Context, n *models.Note) (int, error) {
	db, err := connectDb()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// TODO: Необходимо реализовать запрет на создание одинаковых заметок по name
	q := `INSERT INTO notes (person_id, category_id, name) VALUES ($1, $2, $3) RETURNING id`
	var id int
	if err := db.QueryRowContext(ctx, q, n.PersonId, n.CategoryId, n.Name).Scan(&id); err != nil {
		log.Errorf("cannot save note: %v", err)
		return 0, err
	}

	return id, nil
}

func GetNote(ctx context.Context, id int) (*models.Note, error) {
	db, err := connectDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var n models.Note
	q := `SELECT * FROM notes WHERE id = $1`

	if err := db.QueryRowContext(ctx, q, id).Scan(&n.Id, &n.PersonId, &n.CategoryId, &n.Name, &n.Created); err != nil {
		log.Errorf("err: %v, id: %d", err, id)
		return nil, err
	}
	return &n, nil
}

func GetAllNotes(ctx context.Context) (*[]models.Note, error) {
	db, err := connectDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	notes := make([]models.Note, 0)
	q := `SELECT * FROM notes ORDER BY id ASC`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	defer rows.Close()

	var note models.Note
	for rows.Next() {
		if err := rows.Scan(&note.Id, &note.PersonId, &note.CategoryId, &note.Name, &note.Created); err != nil {
			log.Errorf("err in rows scan: %v", err)
			continue
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		log.Errorf("err in rows.Err: %v", err)
	}
	return &notes, nil
}

func DeleteNote(ctx context.Context, id int) error {
	db, err := connectDb()
	if err != nil {
		return err
	}
	defer db.Close()

	// Проверяем есть ли заметка с таким айди
	var exist bool
	q := `SELECT id FROM notes WHERE id = $1`

	if err := db.QueryRowContext(ctx, q, id).Scan(&exist); err == sql.ErrNoRows {
		log.Errorf("err: %v, id: %d", err, id)
		return err
	}

	// Если заметка с переданным айди есть - удаляем
	q = `DELETE FROM notes WHERE id = $1`

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		log.Errorf("err: %v, id: %d", err, id)
		return err
	}

	return nil
}
