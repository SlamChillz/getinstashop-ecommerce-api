package tests

import (
	"github.com/google/uuid"
	mockdb "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/mock"
	"github.com/slamchillz/getinstashop-ecommerce-api/pkg/token"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListAllProduct(t *testing.T) {
	//_, hashPass := randomPassword(t)
	userId := uuid.New()
	testCases := []struct {
		name     string
		auth     func(t *testing.T, req *http.Request, tokenCreator *token.JWT)
		stubs    func(store *mockdb.MockStore)
		response func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Admin",
			auth: func(t *testing.T, req *http.Request, tokenCreator *token.JWT) {
				addAuthorization(t, req, tokenCreator, userId, true)
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllProduct(gomock.Any()).
					Times(1)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Authenticated User",
			auth: func(t *testing.T, req *http.Request, tokenCreator *token.JWT) {
				addAuthorization(t, req, tokenCreator, userId, false)
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllProduct(gomock.Any()).
					Times(1)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Unauthenticated User",
			auth: func(t *testing.T, req *http.Request, tokenCreator *token.JWT) {
				//addAuthorization(t, req, tokenCreator, userId, true)
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllProduct(gomock.Any()).
					Times(0)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.stubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/api/v1/products"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.auth(t, request, server.TokenCreator())
			server.Router().ServeHTTP(recorder, request)
			tc.response(t, recorder)
		})
	}
}
