package file

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/cchrisris/go-course/hw3/internal/domain"
	"github.com/cchrisris/go-course/hw3/internal/ports"
)

type Repository struct {
	mutex   sync.Mutex
	path    string
	data    map[string]int64
}

func NewRepository(path string) (*Repository, error) {
	repo := &Repository{
		path: path,
		data: make(map[string]int64),
	}
	if err := repo.loadFromDisk(); err != nil {
		return nil, err
	}
	return repo, nil
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
	if err := fn(transactionStore); err != nil {
		return err
	}
	return repo.flushToDisk()
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

func (repo *Repository) loadFromDisk() error {
	if repo.path == "" {
		return errors.New("file repository path is empty")
	}
	if err := os.MkdirAll(filepath.Dir(repo.path), 0o755); err != nil {
		return err
	}

	f, err := os.Open(repo.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	var rawMap map[string]int64
	if err := json.NewDecoder(f).Decode(&rawMap); err != nil {
		return err
	}
	for accountID, balance := range rawMap {
		repo.data[accountID] = balance
	}
	return nil
}

func (repo *Repository) flushToDisk() error {
	tmp := repo.path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(repo.data); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, repo.path)
}


