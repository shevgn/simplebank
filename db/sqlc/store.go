package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store is a database store interface
type Store interface {
	Querier
	TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error)
}

// SQLStore is a database store
type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

// NewStore creates a new Store
func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			return fmt.Errorf("tx error: %w, rollback error: %w", err, rbErr)
		}
	}

	return tx.Commit(ctx)
}

// TransferTxParams is a set of parameters for TransferTx
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is a result of TransferTx
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *SQLStore) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(args))
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		if args.FromAccountID < args.ToAccountID {
			result.FromAccount, result.ToAccount, err = addBalance(
				ctx,
				q,
				args.FromAccountID,
				-args.Amount,
				args.ToAccountID,
				args.Amount,
			)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addBalance(
				ctx,
				q,
				args.ToAccountID,
				args.Amount,
				args.FromAccountID,
				-args.Amount,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func addBalance(
	ctx context.Context,
	q *Queries,
	accountFromID int64,
	amountFrom int64,
	accountToID int64,
	amountTo int64,
) (accountFrom Account, accountTo Account, err error) {
	accountFrom, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountFromID,
		Amount: amountFrom,
	})
	if err != nil {
		return
	}

	accountTo, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountToID,
		Amount: amountTo,
	})
	if err != nil {
		return
	}

	return
}
