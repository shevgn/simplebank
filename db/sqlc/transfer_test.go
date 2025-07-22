package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	arg := CreateTransferParams{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        100,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func transfersEqual(t *testing.T, a, b Transfer) {
	t.Helper()

	require.Equal(t, a.ID, b.ID)
	require.Equal(t, a.FromAccountID, b.FromAccountID)
	require.Equal(t, a.ToAccountID, b.ToAccountID)
	require.Equal(t, a.Amount, b.Amount)

	require.WithinDuration(t, a.CreatedAt.Time, b.CreatedAt.Time, time.Second)
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transferNew := createRandomTransfer(t)

	transfer, err := testQueries.GetTransfer(context.Background(), transferNew.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	transfersEqual(t, transferNew, transfer)
}

func TestListTransfers(t *testing.T) {
	transfersCount := 10
	listTransfersLimit := 5
	listTransfersOffset := 5

	for range transfersCount {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		FromAccountID: 1,
		ToAccountID:   2,
		Limit:         int32(listTransfersLimit),
		Offset:        int32(listTransfersOffset),
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, listTransfersLimit)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
