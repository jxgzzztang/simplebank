package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jxgzzztang/simplebank/db/mock"
	db "github.com/jxgzzztang/simplebank/db/sqlc"
	"github.com/jxgzzztang/simplebank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAccount(t *testing.T)  {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	user, _ := RandomUser(t)
	account := randomAccount(user)

	testCases := []struct {
		Name string
		AccountID int64
		BuildStubs func(store *mock.MockStore)
		SetupAuth func(t *testing.T, request *http.Request)
		CheckResponse func(response *httptest.ResponseRecorder)
	}{
		{
			Name: "OK",
			AccountID: account.ID,
			SetupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, user.Username, time.Minute)
			},
			BuildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			CheckResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireMatchRequestBody(t, recorder.Body, account)
			},
		},
		{
			Name: "UnauthorizedUser",
			AccountID: account.ID,
			SetupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, "unauthorized_user", time.Minute)
			},
			BuildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			CheckResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			Name: "NotFound",
			AccountID: account.ID,
			SetupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, user.Username, time.Minute)
			},
			BuildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, pgx.ErrNoRows)
			},
			CheckResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			Name: "InternalServer",
			AccountID: account.ID,
			SetupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, user.Username, time.Minute)
			},
			BuildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, pgx.ErrTxClosed)
			},
			CheckResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			Name: "BadRequest",
			AccountID: 0,
			SetupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, user.Username, time.Minute)
			},
			BuildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			CheckResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			store := mock.NewMockStore(ctrl)
			tc.BuildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/account/%d", tc.AccountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			tc.SetupAuth(t, req)
			server.router.ServeHTTP(recorder, req)
			tc.CheckResponse(recorder)

		})
	}
}

func randomAccount(user db.User) db.Account {
	return db.Account{
		ID: util.RandomInt(1, 199),
		Owner: user.Username,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireMatchRequestBody(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var actual db.Account

	err = json.Unmarshal(data, &actual)
	require.NoError(t, err)
	require.Equal(t, account, actual)

}