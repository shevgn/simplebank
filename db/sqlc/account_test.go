package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shevgn/simplebank/util"
	"github.com/stretchr/testify/require"
)

// createRandomAccount creates random account for testing
func createRandomAccount(t *testing.T) Account {

	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func accountsEqual(t *testing.T, a, b Account) {
	t.Helper()

	require.Equal(t, a.ID, b.ID)
	require.Equal(t, a.Owner, b.Owner)
	require.Equal(t, a.Balance, b.Balance)
	require.Equal(t, a.Currency, b.Currency)

	require.WithinDuration(t, a.CreatedAt.Time, b.CreatedAt.Time, time.Second)
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	accountNew := createRandomAccount(t)

	account, err := testQueries.GetAccount(context.Background(), accountNew.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	accountsEqual(t, accountNew, account)
}

func TestGetAccountForUpdate(t *testing.T) {
	accountNew := createRandomAccount(t)

	account, err := testQueries.GetAccountForUpdate(context.Background(), accountNew.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	accountsEqual(t, accountNew, account)
}

func TestUpdateAccount(t *testing.T) {
	accountNew := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      accountNew.ID,
		Balance: util.RandomBalance(),
	}

	account, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	accountNew.Balance = arg.Balance

	accountsEqual(t, accountNew, account)
}

func TestDeleteAccount(t *testing.T) {
	accountNew := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), accountNew.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), accountNew.ID)

	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestListAccounts(t *testing.T) {
	accountsCount := 10
	listAccountsLimit := 5
	listAccountsOffset := 5

	for range accountsCount {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  int32(listAccountsLimit),
		Offset: int32(listAccountsOffset),
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, listAccountsLimit)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestAddAccountBalance(t *testing.T) {
	accountNew := createRandomAccount(t)

	arg := AddAccountBalanceParams{
		ID:     accountNew.ID,
		Amount: util.RandomBalance(),
	}

	account, err := testQueries.AddAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	accountNew.Balance += arg.Amount

	accountsEqual(t, accountNew, account)
}
