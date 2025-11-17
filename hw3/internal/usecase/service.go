package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/cchrisris/go-course/hw3/internal/domain"
	"github.com/cchrisris/go-course/hw3/internal/ports"
)

var (
	ErrInvalidAmount      = errors.New("amount must be positive")
	ErrSameAccount        = errors.New("cannot transfer to the same account")
	ErrInsufficientFunds  = errors.New("insufficient funds")
)

type Service interface {
	Deposit(ctx context.Context, accountID string, amount int64) (int64, error)
	Transfer(ctx context.Context, fromID, toID string, amount int64) error
	GetBalance(ctx context.Context, accountID string) (int64, error)
}

func NewService(repo ports.Repository) Service {
	return &service{repo: repo}
}

type service struct {
	repo ports.Repository
}

func (s *service) Deposit(ctx context.Context, accountID string, amount int64) (int64, error) {
	if err := validateAccountID(accountID); err != nil {
		return 0, err
	}
	if amount <= 0 {
		return 0, ErrInvalidAmount
	}

	var newBalance int64
	err := s.repo.RunInTx(ctx, func(store ports.ReadWriter) error {
		account, err := store.Load(ctx, accountID)
		switch {
		case err == nil:
			// ok
		case errors.Is(err, ports.ErrNotFound):
			account = domain.Account{ID: accountID, Balance: 0}
		default:
			return err
		}

		account.Balance += amount
		if err := store.Save(ctx, account); err != nil {
			return err
		}
		newBalance = account.Balance
		return nil
	})
	if err != nil {
		return 0, err
	}
	return newBalance, nil
}

func (s *service) Transfer(ctx context.Context, fromID, toID string, amount int64) error {
	if err := validateAccountID(fromID); err != nil {
		return err
	}
	if err := validateAccountID(toID); err != nil {
		return err
	}
	if fromID == toID {
		return ErrSameAccount
	}
	if amount <= 0 {
		return ErrInvalidAmount
	}

	return s.repo.RunInTx(ctx, func(store ports.ReadWriter) error {
		from, err := store.Load(ctx, fromID)
		if err != nil {
			return err
		}
		to, err := store.Load(ctx, toID)
		if err != nil {
			return err
		}
		if from.Balance < amount {
			return ErrInsufficientFunds
		}

		from.Balance -= amount
		to.Balance += amount
		if err := store.Save(ctx, from); err != nil {
			return err
		}
		if err := store.Save(ctx, to); err != nil {
			return err
		}
		return nil
	})
}

func (s *service) GetBalance(ctx context.Context, accountID string) (int64, error) {
	if err := validateAccountID(accountID); err != nil {
		return 0, err
	}
	account, err := s.repo.Load(ctx, accountID)
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}

func validateAccountID(id string) error {
	if len(id) == 0 {
		return fmt.Errorf("account id is required")
	}
	return nil
}


