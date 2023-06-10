package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	mockdb "github.com/aalug/go-bank/db/mock"
	db "github.com/aalug/go-bank/db/sqlc"
	"github.com/aalug/go-bank/token"
	"github.com/aalug/go-bank/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateTransferAPI(t *testing.T) {
	var amount int64 = 10

	user1, _ := generateRandomUser(t)
	user2, _ := generateRandomUser(t)
	user3, _ := generateRandomUser(t)

	// 1 and 2 with the same currency, 3 with  different currency
	account1eur := generateRandomAccount(user1.Username)
	account2eur := generateRandomAccount(user2.Username)
	account3usd := generateRandomAccount(user3.Username)

	account1eur.Currency = utils.EUR
	account2eur.Currency = utils.EUR
	account3usd.Currency = utils.USD

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, r *http.Request, maker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": account1eur.ID,
				"to_account_id":   account2eur.ID,
				"amount":          amount,
				"currency":        utils.EUR,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// get accounts
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1eur.ID)).
					Times(1).
					Return(account1eur, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2eur.ID)).
					Times(1).
					Return(account2eur, nil)

				params := db.TransferTxParams{
					FromAccountID: account1eur.ID,
					ToAccountID:   account2eur.ID,
					Amount:        amount,
				}

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(params)).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "From Account Not Found",
			body: gin.H{
				"from_account_id": account1eur.ID,
				"to_account_id":   account2eur.ID,
				"amount":          amount,
				"currency":        utils.EUR,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1eur.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2eur.ID)).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "To Account Not Found",
			body: gin.H{
				"from_account_id": account1eur.ID,
				"to_account_id":   account2eur.ID,
				"amount":          amount,
				"currency":        utils.EUR,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1eur.ID)).
					Times(1).
					Return(account1eur, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2eur.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "From Account Currency Mismatch",
			body: gin.H{
				"from_account_id": account3usd.ID,
				"to_account_id":   account2eur.ID,
				"amount":          amount,
				"currency":        utils.EUR,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, user3.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account3usd.ID)).
					Times(1).
					Return(account3usd, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2eur.ID)).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "To Account Currency Mismatch",
			body: gin.H{
				"from_account_id": account1eur.ID,
				"to_account_id":   account3usd.ID,
				"amount":          amount,
				"currency":        utils.EUR,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1eur.ID)).
					Times(1).
					Return(account1eur, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account3usd.ID)).
					Times(1).
					Return(account3usd, nil)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Currency",
			body: gin.H{
				"from_account_id": account1eur.ID,
				"to_account_id":   account2eur.ID,
				"amount":          amount,
				"currency":        "INVALID",
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Negative Amount",
			body: gin.H{
				"from_account_id": account1eur.ID,
				"to_account_id":   account2eur.ID,
				"amount":          -amount,
				"currency":        utils.EUR,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Amount",
			body: gin.H{
				"from_account_id": account1eur.ID,
				"to_account_id":   account2eur.ID,
				"amount":          10.50,
				"currency":        utils.EUR,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "GetAccount Error",
			body: gin.H{
				"from_account_id": account1eur.ID,
				"to_account_id":   account2eur.ID,
				"amount":          amount,
				"currency":        utils.EUR,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "TransferTx Error",
			body: gin.H{
				"from_account_id": account1eur.ID,
				"to_account_id":   account2eur.ID,
				"amount":          amount,
				"currency":        utils.EUR,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1eur.ID)).
					Times(1).
					Return(account1eur, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2eur.ID)).
					Times(1).
					Return(account2eur, nil)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.TransferTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/transfers"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}
