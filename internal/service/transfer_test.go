package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/example/ledger/internal/db"
	"github.com/example/ledger/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestTransferTxDeadlock(t *testing.T) {

	connString := "postgresql://root:secret@localhost:5432/ledger?sslmode=disable"
	connPool, err := db.NewConnectionPool(context.Background(), connString)
	if err != nil {
		t.Skip("Skipping test: cannot connect to db:", err)
	}
	defer connPool.Close()

	store := repository.NewStore(connPool)
	svc := NewService(store)

	// Create accounts
	account1, err := svc.CreateAccount(context.Background(), 1000, "USD")
	require.NoError(t, err)
	account2, err := svc.CreateAccount(context.Background(), 1000, "USD")
	require.NoError(t, err)

	n := 10
	amount := int64(10)
	errs := make(chan error)
	results := make(chan repository.TransferTxResult)

	// Run n concurrent transfer transactions
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), "txName", txName)
			result, err := svc.Transfer(ctx, account1.ID, account2.ID, amount)
			errs <- err
			results <- result
		}()
	}

	// Check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)
		require.Equal(t, account1.ID, result.FromAccount.ID)
		require.Equal(t, account2.ID, result.ToAccount.ID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)
	}

	// Check final updated balance
	updatedAccount1, err := svc.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := svc.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
