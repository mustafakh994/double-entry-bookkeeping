package service

import (
	"context"

	"github.com/example/ledger/internal/repository"
)

type Service struct {
	store repository.Store
}

func NewService(store repository.Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateAccount(ctx context.Context, balance int64, currency string) (repository.Account, error) {
	return s.store.CreateAccount(ctx, repository.CreateAccountParams{
		Balance:  balance,
		Currency: currency,
	})
}

func (s *Service) GetAccount(ctx context.Context, id int64) (repository.Account, error) {
	return s.store.GetAccount(ctx, id)
}

func (s *Service) Transfer(ctx context.Context, fromID, toID, amount int64) (repository.TransferTxResult, error) {
	return s.store.TransferTx(ctx, repository.TransferTxParams{
		FromAccountID: fromID,
		ToAccountID:   toID,
		Amount:        amount,
	})
}
