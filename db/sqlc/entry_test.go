package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	arg := CreateEntryParams{
		AccountID: 1,
		Amount:    100,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func entriesEqual(t *testing.T, a, b Entry) {
	t.Helper()

	require.Equal(t, a.ID, b.ID)
	require.Equal(t, a.AccountID, b.AccountID)
	require.Equal(t, a.Amount, b.Amount)

	require.WithinDuration(t, a.CreatedAt.Time, b.CreatedAt.Time, time.Second)
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entryNew := createRandomEntry(t)

	entry, err := testQueries.GetEntry(context.Background(), entryNew.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	entriesEqual(t, entryNew, entry)
}

func TestListEntries(t *testing.T) {
	entriesCount := 10
	listEntriesLimit := 5
	listEntriesOffset := 5

	for range entriesCount {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		AccountID: 1,
		Limit:     int32(listEntriesLimit),
		Offset:    int32(listEntriesOffset),
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, listEntriesLimit)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
