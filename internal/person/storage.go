package person

import (
	"context"
	"github.com/beloslav13/servernotes/internal/models"
)

type Repository interface {
	Create(ctx context.Context, n *models.Person) (int, error)
	Get(ctx context.Context, id int) (*models.Person, error)
	GetAll(ctx context.Context) (*[]models.Person, error)
	Delete(ctx context.Context, id int) error
}
