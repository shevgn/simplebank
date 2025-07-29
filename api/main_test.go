package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/shevgn/simplebank/db/sqlc"
	"github.com/shevgn/simplebank/util"
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := &util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	return NewServer(config, store)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
