package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx"
	mockdb "github.com/shevgn/simplebank/db/mock"
	db "github.com/shevgn/simplebank/db/sqlc"
	"github.com/shevgn/simplebank/token"
	"github.com/shevgn/simplebank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTransferAPI(t *testing.T) {
	amount := util.RandomInt(1, 100)

	userFrom, _ := randomUser(t)
	userTo, _ := randomUser(t)
	userWithCurrencyMismatch, _ := randomUser(t)

	accountFrom := randomAccount(userFrom.Username)
	accountTo := randomAccount(userTo.Username)
	accountWithCurrencyMismatch := randomAccount(userWithCurrencyMismatch.Username)

	accountFrom.Currency = util.USD
	accountTo.Currency = util.USD
	accountWithCurrencyMismatch.Currency = util.EUR

	testCases := []struct {
		name          string
		requestBody   CreateTransferRequest
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			requestBody: CreateTransferRequest{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
				Currency:      util.USD,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, userFrom.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(accountFrom.ID)).
					Times(1).
					Return(accountFrom, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(accountTo.ID)).
					Times(1).
					Return(accountTo, nil)

				arg := db.TransferTxParams{
					FromAccountID: accountFrom.ID,
					ToAccountID:   accountTo.ID,
					Amount:        amount,
				}

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			requestBody: CreateTransferRequest{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
				Currency:      util.USD,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, userFrom.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(accountFrom.ID)).
					Times(1).
					Return(db.Account{}, pgx.ErrNoRows)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(accountTo.ID)).
					Times(0)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			requestBody: CreateTransferRequest{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
				Currency:      util.USD,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, userFrom.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(accountFrom.ID)).
					Times(1).
					Return(accountFrom, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(accountTo.ID)).
					Times(1).
					Return(db.Account{}, pgx.ErrNoRows)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "FromAccountCurrencyMismatch",
			requestBody: CreateTransferRequest{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
				Currency:      util.USD,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, userFrom.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(accountFrom.ID)).
					Times(1).
					Return(accountWithCurrencyMismatch, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(accountTo.ID)).
					Times(0)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToAccountCurrencyMismatch",
			requestBody: CreateTransferRequest{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
				Currency:      util.USD,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, userFrom.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(accountFrom.ID)).
					Times(1).
					Return(accountFrom, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(accountTo.ID)).
					Times(1).
					Return(accountWithCurrencyMismatch, nil)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			requestBody: CreateTransferRequest{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
				Currency:      "invalid",
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, userFrom.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidAmount",
			requestBody: CreateTransferRequest{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        -amount,
				Currency:      util.USD,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, userFrom.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "GetAccountError",
			requestBody: CreateTransferRequest{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
				Currency:      util.USD,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, userFrom.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(accountFrom.ID)).
					Times(1).
					Return(db.Account{}, pgx.ErrClosedPool)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(accountTo.ID)).
					Times(0)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "TransferTxError",
			requestBody: CreateTransferRequest{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
				Currency:      util.USD,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, userFrom.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(accountFrom.ID)).
					Times(1).
					Return(accountFrom, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(accountTo.ID)).
					Times(1).
					Return(accountTo, nil)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.TransferTxResult{}, pgx.ErrClosedPool)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/transfers"
			body, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
