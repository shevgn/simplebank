package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)
	fmt.Println(">> balance before transfer", accountFrom.Balance, accountTo.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)

	for range n {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result

		}()
	}

	exists := make(map[int]bool)
	for range n {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)
		require.NotEmpty(t, result.Transfer)
		require.NotEmpty(t, result.FromAccount)
		require.NotEmpty(t, result.ToAccount)
		require.NotEmpty(t, result.FromEntry)
		require.NotEmpty(t, result.ToEntry)

		transfer := result.Transfer
		require.Equal(t, accountFrom.ID, transfer.FromAccountID)
		require.Equal(t, accountTo.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, accountFrom.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, accountTo.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.Equal(t, accountFrom.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.Equal(t, accountTo.ID, toAccount.ID)

		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diffFrom := accountFrom.Balance - fromAccount.Balance
		diffTo := toAccount.Balance - accountTo.Balance
		require.Equal(t, diffFrom, diffTo)
		require.True(t, diffFrom > 0)
		require.True(t, diffTo%amount == 0)

		k := int(diffTo / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, exists, k)
		exists[k] = true
	}

	updatedFromAccount, err := store.GetAccount(context.Background(), accountFrom.ID)
	require.NoError(t, err)

	updatedToAccount, err := store.GetAccount(context.Background(), accountTo.ID)
	require.NoError(t, err)

	fmt.Println(">> balance after transfers", updatedFromAccount.Balance, updatedToAccount.Balance)

	require.Equal(t, accountFrom.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, accountTo.Balance+int64(n)*amount, updatedToAccount.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)
	fmt.Println(">> balance before transfer", accountFrom.Balance, accountTo.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error, n)

	for i := range n {
		fromAccountID := accountFrom.ID
		toAccountID := accountTo.ID

		if i%2 == 1 {
			fromAccountID = accountTo.ID
			toAccountID = accountFrom.ID
		}

		go func() {
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for range n {
		err := <-errs
		require.NoError(t, err)
	}

	updatedFromAccount, err := store.GetAccount(context.Background(), accountFrom.ID)
	require.NoError(t, err)

	updatedToAccount, err := store.GetAccount(context.Background(), accountTo.ID)
	require.NoError(t, err)

	fmt.Println(">> balance after transfers", updatedFromAccount.Balance, updatedToAccount.Balance)

	require.Equal(t, accountFrom.Balance, updatedFromAccount.Balance)
	require.Equal(t, accountTo.Balance, updatedToAccount.Balance)
}

// func TestStore_TransferTx(t *testing.T) {
// 	tests := []struct {
// 		name string // description of this test case
// 		// Named input parameters for receiver constructor.
// 		db *pgx.Conn
// 		// Named input parameters for target function.
// 		args    db.TransferTxParams
// 		want    db.TransferTxResult
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := db.NewStore(tt.db)
// 			got, gotErr := s.TransferTx(context.Background(), tt.args)
// 			if gotErr != nil {
// 				if !tt.wantErr {
// 					t.Errorf("TransferTx() failed: %v", gotErr)
// 				}
// 				return
// 			}
// 			if tt.wantErr {
// 				t.Fatal("TransferTx() succeeded unexpectedly")
// 			}
// 			// TODO: update the condition below to compare got with tt.want.
// 			if true {
// 				t.Errorf("TransferTx() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
