package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shevgn/simplebank/token"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	req *http.Request,
	tm token.Maker,
	authType string,
	username string,
	duration time.Duration,
) {
	token, _, err := tm.CreateToken(username, duration)
	require.NoError(t, err)

	req.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authType, token))
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name          string
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		chechResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			chechResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorizationHeader",
			setupAuth: func(_ *testing.T, _ *http.Request, _ token.Maker) {
			},
			chechResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorizationType",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, "unsupported", "username", time.Minute)
			},
			chechResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, "", "username", time.Minute)
			},
			chechResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, req *http.Request, tm token.Maker) {
				addAuthorization(t, req, tm, authorizationTypeBearer, "username", -time.Minute)
			},
			chechResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(authPath, authMiddleware(server.tokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{
					"message": "OK",
				})
			})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tt.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tt.chechResponse(t, recorder)
		})
	}
}
