package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/slamchillz/getinstashop-ecommerce-api/cmd/server"
	"github.com/slamchillz/getinstashop-ecommerce-api/config"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func newTestServer(t *testing.T, store db.Store) *server.Server {
	cfg, err := config.LoadConfig("../")
	require.NoError(t, err)
	apiServer, err := server.NewServer(cfg, store)
	require.NoError(t, err)
	return apiServer
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
