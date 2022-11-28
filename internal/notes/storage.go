package notes

import (
	"context"
	"github.com/beloslav13/servernotes/internal/models"
)

type Repository interface {
	Create(ctx context.Context, n *models.Note) (int, error)
	Get(ctx context.Context, id int) (*models.Note, error)
	GetAll(ctx context.Context) (*[]models.Note, error)
	Delete(ctx context.Context, id int) error
}
