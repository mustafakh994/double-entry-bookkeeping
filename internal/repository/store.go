package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	GetAccount(ctx context.Context, id int64) (Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (Account, error)
	UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceParams) error
	CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error)
	ListTransactions(ctx context.Context, arg ListTransactionsParams) ([]Transaction, error)
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transaction `json:"transfer"`
	FromAccount Account     `json:"from_account"`
	ToAccount   Account     `json:"to_account"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 1. Create Transaction Record
		result.Transfer, err = q.CreateTransaction(ctx, CreateTransactionParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// 2. Update Account Balances with Locking
		// To avoid deadlocks, we update accounts in a consistent order (by ID)
		// However, the user requirement specifically asks for "Select For Update"
		// We need to ensure we lock them in order to avoid deadlocks if we were doing arbitrary updates,
		// but for a simple transfer, sorting by ID is the standard way to avoid deadlocks.

		var account1, account2 Account
		
		// Lock the accounts
		if arg.FromAccountID < arg.ToAccountID {
			account1, err = q.GetAccountForUpdate(ctx, arg.FromAccountID)
			if err != nil { return err }
			account2, err = q.GetAccountForUpdate(ctx, arg.ToAccountID)
			if err != nil { return err }
		} else {
			account2, err = q.GetAccountForUpdate(ctx, arg.ToAccountID)
			if err != nil { return err }
			account1, err = q.GetAccountForUpdate(ctx, arg.FromAccountID)
			if err != nil { return err }
		}

		// Check balance
		if account1.ID == arg.FromAccountID {
			// account1 is the sender
			if account1.Balance < arg.Amount {
				return fmt.Errorf("insufficient funds")
			}
			// Update balances
			err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
				ID:      arg.FromAccountID,
				Balance: account1.Balance - arg.Amount,
			})
			if err != nil { return err }
			
			err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
				ID:      arg.ToAccountID,
				Balance: account2.Balance + arg.Amount,
			})
			if err != nil { return err }

			result.FromAccount = account1
			result.FromAccount.Balance -= arg.Amount
			result.ToAccount = account2
			result.ToAccount.Balance += arg.Amount
		} else {
			// account2 is the sender (this case happens if ToID < FromID)
			if account2.Balance < arg.Amount {
				return fmt.Errorf("insufficient funds")
			}
			// Update balances
			err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
				ID:      arg.FromAccountID,
				Balance: account1.Balance + arg.Amount,
			})
			if err != nil { return err }

			err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
				ID:      arg.ToAccountID,
				Balance: account2.Balance - arg.Amount,
			})
			if err != nil { return err }

			result.FromAccount = account2
			result.FromAccount.Balance -= arg.Amount
			result.ToAccount = account1
			result.ToAccount.Balance += arg.Amount
		}

		return nil
	})

	return result, err
}
