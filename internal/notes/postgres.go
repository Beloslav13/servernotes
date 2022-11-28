package notes

import (
	"context"
	"database/sql"
	"github.com/beloslav13/servernotes/internal/models"
	"github.com/beloslav13/servernotes/pkg/logger"
	"github.com/beloslav13/servernotes/pkg/postgresql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type repository struct {
	logger logger.Logger
}

func NewRepository(logger logger.Logger) Repository {
	return &repository{
		logger: logger,
	}
}

func (r *repository) Create(ctx context.Context, n *models.Note) (int, error) {
	db, err := postgresql.ConnectDb(r.logger)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// TODO: Необходимо реализовать запрет на создание одинаковых заметок по name
	q := `INSERT INTO notes (person_id, category_id, name) VALUES ($1, $2, $3) RETURNING id`
	var id int
	if err := db.QueryRowContext(ctx, q, n.PersonId, n.CategoryId, n.Name).Scan(&id); err != nil {
		r.logger.Errorf("cannot save note: %v", err)
		return 0, err
	}

	return id, nil
}

func (r *repository) Get(ctx context.Context, id int) (*models.Note, error) {
	db, err := postgresql.ConnectDb(r.logger)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var n models.Note
	q := `SELECT * FROM notes WHERE id = $1`

	if err := db.QueryRowContext(ctx, q, id).Scan(&n.Id, &n.PersonId, &n.CategoryId, &n.Name, &n.Created); err != nil {
		r.logger.Errorf("err: %v, id: %d", err, id)
		return nil, err
	}
	return &n, nil
}

func (r *repository) GetAll(ctx context.Context) (*[]models.Note, error) {
	db, err := postgresql.ConnectDb(r.logger)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	notes := make([]models.Note, 0)
	q := `SELECT * FROM notes ORDER BY id ASC`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		r.logger.Errorf("err: %v", err)
		return nil, err
	}
	defer rows.Close()

	var note models.Note
	for rows.Next() {
		if err := rows.Scan(&note.Id, &note.PersonId, &note.CategoryId, &note.Name, &note.Created); err != nil {
			r.logger.Errorf("err in rows scan: %v", err)
			continue
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		r.logger.Errorf("err in rows.Err: %v", err)
	}
	return &notes, nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	db, err := postgresql.ConnectDb(r.logger)
	if err != nil {
		return err
	}
	defer db.Close()

	// Проверяем есть ли заметка с таким айди
	var exist bool
	q := `SELECT id FROM notes WHERE id = $1`

	if err := db.QueryRowContext(ctx, q, id).Scan(&exist); err == sql.ErrNoRows {
		r.logger.Errorf("err: %v, id: %d", err, id)
		return err
	}

	// Если заметка с переданным айди есть - удаляем
	q = `DELETE FROM notes WHERE id = $1`

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		r.logger.Errorf("err: %v, id: %d", err, id)
		return err
	}

	return nil
}
