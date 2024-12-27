package tests

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/slamchillz/getinstashop-ecommerce-api/cmd/server"
	"github.com/slamchillz/getinstashop-ecommerce-api/config"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/constants"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/utils"
	"github.com/slamchillz/getinstashop-ecommerce-api/pkg/token"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

var testUserId = uuid.New()

func newTestServer(t *testing.T, store db.Store) *server.Server {
	cfg, err := config.LoadConfig("../")
	require.NoError(t, err)
	apiServer, err := server.NewServer(cfg, store)
	require.NoError(t, err)
	return apiServer
}

func randomPassword(t *testing.T) (string, string) {
	pass := "password123"
	hashPass, err := utils.HashPassword(pass)
	require.NoError(t, err)
	return pass, hashPass
}

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenCreator *token.JWT,
	username uuid.UUID,
	admin bool,
) {
	jwtToken, err := tokenCreator.CreateToken(username, admin)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", constants.AuthenticationScheme, jwtToken)
	request.Header.Set(constants.AuthenticationHeader, authorizationHeader)
}

//func randomProduct() db.Product {
//	return db.Product{
//		ID:          uuid.New(),
//		Name:        "test",
//		Description: "test product",
//		Price:       1000,
//		Stock:       5,
//		CreatedBy:   testUserId,
//	}
//}
