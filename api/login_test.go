package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jxgzzztang/simplebank/db/mock"
	db "github.com/jxgzzztang/simplebank/db/sqlc"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLogin(t *testing.T) {
	user, password := RandomUser(t)

	testCases := []struct {
		Name string
		Body gin.H
		BuildStep func(store *mock.MockStore)
		CheckResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			Name: "ok",
			Body: gin.H{
				"username": user.Username,
				"password": password,
			},
			BuildStep: func(store *mock.MockStore) {
				// 现有的期望调用
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
				
				// 添加对 CreateSession 的期望调用
				store.EXPECT().CreateSessions(gomock.Any(), gomock.Any()).Times(1).Return(db.Session{}, nil)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mock.NewMockStore(ctrl)
			tc.BuildStep(store)
			body, err := json.Marshal(tc.Body)
			require.NoError(t, err)
			server := NewServer(store)

			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

		})
	}

}
