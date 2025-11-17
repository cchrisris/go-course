package memory

import (
	"context"
	"sync"

	"github.com/cchrisris/go-course/hw3/internal/domain"
	"github.com/cchrisris/go-course/hw3/internal/ports"
)

type Repository struct {
	mutex sync.Mutex
	data  map[string]int64
}

func NewRepository() *Repository {
	return &Repository{
		data: make(map[string]int64),
	}
}

func (repo *Repository) Load(ctx context.Context, id string) (domain.Account, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	balance, exists := repo.data[id]
	if !exists {
		return domain.Account{}, ports.ErrNotFound
	}
	return domain.Account{ID: id, Balance: balance}, nil
}

func (repo *Repository) RunInTx(ctx context.Context, fn func(store ports.ReadWriter) error) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	transactionStore := &transactionStore{data: repo.data}
	return fn(transactionStore)
}

type transactionStore struct {
	data map[string]int64
}

func (store *transactionStore) Load(ctx context.Context, id string) (domain.Account, error) {
	balance, exists := store.data[id]
	if !exists {
		return domain.Account{}, ports.ErrNotFound
	}
	return domain.Account{ID: id, Balance: balance}, nil
}

func (store *transactionStore) Save(ctx context.Context, account domain.Account) error {
	store.data[account.ID] = account.Balance
	return nil
}


