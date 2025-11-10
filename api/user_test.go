package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jxgzzztang/simplebank/db/mock"
	db "github.com/jxgzzztang/simplebank/db/sqlc"
	"github.com/jxgzzztang/simplebank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)



func RandomUser(t *testing.T) (user db.User, password string) {
	owner := util.RandomOwner()
	password = util.RandomString(6)
	hashPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	user = db.User{
		Username: owner,
		FullName: util.RandomOwner(),
		Email:    util.RandomEmail(),
		HashedPassword: hashPassword,
	}
	return user, password
}

type CreateUserParamsPwdMatcher struct {
	arg db.CreateUserParams
	password string
}

func (c CreateUserParamsPwdMatcher) Matches (x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(c.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	c.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(c.arg, arg)
}

func (c CreateUserParamsPwdMatcher) String() string {
	return fmt.Sprintf("%v", c.arg)
}

func CheckPasswordMatcher(arg db.CreateUserParams, password string) gomock.Matcher {
	return CreateUserParamsPwdMatcher{
		arg: arg,
		password: password,
	}
}

func TestCreateUser(t *testing.T) {
	user, password := RandomUser(t)

	testCases := []struct {
		Name string
		Body gin.H
		BuildStubs func(store *mock.MockStore)
		CheckResponse func(t *testing.T, response *httptest.ResponseRecorder)
	}{
		{
			Name: "ok",
			Body: gin.H{
				"username": user.Username,
				"full_name": user.FullName,
				"email": user.Email,
				"password": password,
			},
			BuildStubs: func(store *mock.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email: user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(), CheckPasswordMatcher(arg, password)).Times(1).Return(user, nil)
			},
			CheckResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, response.Code)
				requireMatchCreateUserRequestBody(t, response.Body, user)
			},
		},
		{
			Name: "bad request - invalid username",
			Body: gin.H{
				"username": "##(*#&@(*&@#$",
				"full_name": user.FullName,
				"email": user.Email,
				"password": password,
			},
			BuildStubs: func(store *mock.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			CheckResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, response.Code)
			},
		},
		{
			Name: "bad request - invalid email",
			Body: gin.H{
				"username": user.Username,
				"full_name": user.FullName,
				"email": "##(*#&@(*&@#$",
				"password": password,
			},
			BuildStubs: func(store *mock.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			CheckResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, response.Code)
			},
		},
		{
			Name: "bad request - short password",
			Body: gin.H{
				"username": user.Username,
				"full_name": user.FullName,
				"email": user.Email,
				"password": "12345",
			},
			BuildStubs: func(store *mock.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			CheckResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, response.Code)
			},
		},
		{
			Name: "internal server error",
			Body: gin.H{
				"username": user.Username,
				"full_name": user.FullName,
				"email": user.Email,
				"password": password,
			},
			BuildStubs: func(store *mock.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(),gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			CheckResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, response.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockStore := mock.NewMockStore(controller)
			tc.BuildStubs(mockStore)

			server := NewServer(mockStore)
			recorder := httptest.NewRecorder()

			body, err := json.Marshal(tc.Body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/createUser", bytes.NewReader(body))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.CheckResponse(t, recorder)
		})
	}

}

func requireMatchCreateUserRequestBody(t *testing.T, body *bytes.Buffer, user db.User) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)
    
    var actual db.User
    err = json.Unmarshal(data, &actual)
    require.NoError(t, err)
    
    require.Equal(t, user.Username, actual.Username)
    require.Equal(t, user.FullName, actual.FullName)
    require.Equal(t, user.Email, actual.Email)
    require.Empty(t, actual.HashedPassword)
}