package ports

import (
	"context"

	"github.com/cchrisris/go-course/hw3/internal/domain"
)

var ErrNotFound = repositoryNotFoundError("not found")

type repositoryNotFoundError string

func (e repositoryNotFoundError) Error() string { return string(e) }

type ReadWriter interface {
	Load(ctx context.Context, id string) (domain.Account, error)
	Save(ctx context.Context, account domain.Account) error
}

type Repository interface {
	Load(ctx context.Context, id string) (domain.Account, error)

	RunInTx(ctx context.Context, fn func(store ReadWriter) error) error
}
