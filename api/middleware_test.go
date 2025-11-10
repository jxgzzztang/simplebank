package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jxgzzztang/simplebank/db/mock"
	"github.com/jxgzzztang/simplebank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func AddAuthorization(t *testing.T, request *http.Request, username string, duration time.Duration) {
	token, payload, err := util.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)
	request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
}

func TestAuthMiddleware(t *testing.T)  {
	username := util.RandomOwner()

	testCases := []struct {
		Name string
		SetupAuth func(t *testing.T, request *http.Request)
		CheckResponse func(t *testing.T, response *httptest.ResponseRecorder)
	}{
		{
			Name: "ok",
			SetupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, username, time.Minute)
			},
			CheckResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, response.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mock.NewMockStore(ctrl)

			server := NewServer(store)

			authPath := "/auth"
			server.router.GET(authPath, authMiddleware(), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.SetupAuth(t, request)

			server.router.ServeHTTP(recorder, request)

			tc.CheckResponse(t, recorder)
		})
	}

}
