package person

import (
	"context"
	"database/sql"
	"github.com/beloslav13/servernotes/internal/models"
	"github.com/beloslav13/servernotes/pkg/logger"
	"github.com/beloslav13/servernotes/pkg/postgresql"
)

type repository struct {
	logger logger.Logger
}

func NewRepository(logger logger.Logger) Repository {
	return &repository{
		logger: logger,
	}
}

func (r *repository) Create(ctx context.Context, p *models.Person) (int, error) {
	db, err := postgresql.ConnectDb(r.logger)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	q := `
		INSERT INTO
		    persons (tg_user_id, Username)
		VALUES
		    ($1, $2)
		RETURNING id
		`
	r.logger.Trace(`INSERT INTO persons (tg_user_id, Username) VALUES ($1, $2) RETURNING id`)

	var id int
	if err := db.QueryRowContext(ctx, q, p.TgUserId, p.Username).Scan(&id); err != nil {
		r.logger.Errorf("cannot save person: %v", err)
		return 0, err
	}

	return id, nil
}

func (r *repository) Get(ctx context.Context, id int) (*models.Person, error) {
	db, err := postgresql.ConnectDb(r.logger)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var p models.Person
	q := `
		SELECT
		    id, tg_user_id, username, created
		FROM 
		    persons
		WHERE 
		    id = $1
		`
	r.logger.Trace(`SELECT id, tg_user_id, username, created FROM persons WHERE id = $1`)

	if err := db.QueryRowContext(ctx, q, id).Scan(&p.Id, &p.TgUserId, &p.Username, &p.Created); err != nil {
		r.logger.Errorf("err: %v, id: %d", err, id)
		return nil, err
	}
	return &p, nil
}

func (r *repository) GetAll(ctx context.Context) (*[]models.Person, error) {
	db, err := postgresql.ConnectDb(r.logger)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	persons := make([]models.Person, 0)
	q := `
		SELECT
		    id, tg_user_id, username, created
		FROM 
		    persons
		ORDER BY id ASC
		`
	r.logger.Trace(`SELECT id, tg_user_id, username, created FROM persons ORDER BY id ASC `)

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		r.logger.Errorf("err: %v", err)
		return nil, err
	}
	defer rows.Close()

	var person models.Person
	for rows.Next() {
		if err := rows.Scan(&person.Id, &person.TgUserId, &person.Username, &person.Created); err != nil {
			r.logger.Errorf("err in rows scan: %v", err)
			continue
		}
		persons = append(persons, person)
	}
	if err := rows.Err(); err != nil {
		r.logger.Errorf("err in rows.Err: %v", err)
	}
	return &persons, nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	db, err := postgresql.ConnectDb(r.logger)
	if err != nil {
		return err
	}
	defer db.Close()

	// Проверяем есть ли персона с таким айди
	var exist bool
	q := `
		SELECT
		    id
		FROM
		    persons
		WHERE id = $1
		`
	r.logger.Trace(`SELECT id FROM persons WHERE id = $1`)

	if err := db.QueryRowContext(ctx, q, id).Scan(&exist); err == sql.ErrNoRows {
		r.logger.Errorf("err: %v, id: %d", err, id)
		return err
	}

	// Если персона с переданным айди есть - удаляем
	q = `
		DELETE FROM
           persons
		WHERE id = $1
	   `
	r.logger.Trace(`DELETE FROM persons WHERE id = $1`)

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		r.logger.Errorf("err: %v, id: %d", err, id)
		return err
	}

	return nil
}
