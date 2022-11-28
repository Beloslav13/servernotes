package person

import (
	"context"
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
	var id int
	if err := db.QueryRowContext(ctx, q, p.TgUserId, p.Username).Scan(&id); err != nil {
		r.logger.Errorf("cannot save person: %v", err)
		return 0, err
	}

	return id, nil
}

func (r *repository) Get(ctx context.Context, id int) (*models.Person, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) GetAll(ctx context.Context) (*[]models.Person, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
