package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	mockdb "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/mock"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/pkg/token"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	testCases := []struct {
		name     string
		body     gin.H
		auth     func(t *testing.T, req *http.Request, tokenCreator *token.JWT)
		stubs    func(store *mockdb.MockStore)
		response func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: gin.H{
				"description": "IOS device",
				"name":        "iPhone7",
				"price":       120000,
				"stock":       13,
			},
			auth: func(t *testing.T, req *http.Request, tokenCreator *token.JWT) {
				addAuthorization(t, req, tokenCreator, testUserId, true)
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(1)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "Fail",
			body: gin.H{
				"description": "IOS device",
				"name":        "iPhone7",
				//"price":       120000,
				"stock": 13,
			},
			auth: func(t *testing.T, req *http.Request, tokenCreator *token.JWT) {
				addAuthorization(t, req, tokenCreator, testUserId, true)
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Forbidden",
			body: gin.H{
				"description": "IOS device",
				"name":        "iPhone7",
				"price":       120000,
				"stock":       13,
			},
			auth: func(t *testing.T, req *http.Request, tokenCreator *token.JWT) {
				addAuthorization(t, req, tokenCreator, testUserId, false)
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
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
			reqBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/v1/admin/products"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
			require.NoError(t, err)

			tc.auth(t, request, server.TokenCreator())
			server.Router().ServeHTTP(recorder, request)
			tc.response(t, recorder)
		})
	}
}

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

func TestGetSingleProduct(t *testing.T) {
	product := db.GetOneProductRow{
		ID:          uuid.New(),
		Name:        "test",
		Description: "test product",
		Price:       1000,
		Stock:       5,
		CreatedBy:   testUserId,
		CreatedAt:   pgtype.Timestamp{},
		UpdatedAt:   pgtype.Timestamp{},
	}
	testCases := []struct {
		name      string
		productId uuid.UUID
		auth      func(t *testing.T, req *http.Request, tokenCreator *token.JWT)
		stubs     func(store *mockdb.MockStore)
		response  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "Admin",
			productId: product.ID,
			auth: func(t *testing.T, req *http.Request, tokenCreator *token.JWT) {
				addAuthorization(t, req, tokenCreator, testUserId, true)
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetOneProduct(gomock.Any(), gomock.Eq(product.ID)).
					//Return(product, nil).
					Times(1)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "Forbidden",
			productId: product.ID,
			auth: func(t *testing.T, req *http.Request, tokenCreator *token.JWT) {
				addAuthorization(t, req, tokenCreator, testUserId, false)
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetOneProduct(gomock.Any(), gomock.Eq(product.ID)).
					//Return(product, nil).
					Times(0)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
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

			url := fmt.Sprintf("/api/v1/admin/products/%s", tc.productId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.auth(t, request, server.TokenCreator())
			server.Router().ServeHTTP(recorder, request)
			tc.response(t, recorder)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	product := db.GetOneProductRow{
		ID:          uuid.New(),
		Name:        "test",
		Description: "test product",
		Price:       1000,
		Stock:       5,
		CreatedBy:   testUserId,
		CreatedAt:   pgtype.Timestamp{},
		UpdatedAt:   pgtype.Timestamp{},
	}
	testCases := []struct {
		name      string
		productId uuid.UUID
		auth      func(t *testing.T, req *http.Request, tokenCreator *token.JWT)
		stubs     func(store *mockdb.MockStore)
		response  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "Admin",
			productId: product.ID,
			auth: func(t *testing.T, req *http.Request, tokenCreator *token.JWT) {
				addAuthorization(t, req, tokenCreator, testUserId, true)
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteOneProduct(gomock.Any(), gomock.Eq(product.ID)).
					//Return(product, nil).
					Times(1)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name:      "Forbidden",
			productId: product.ID,
			auth: func(t *testing.T, req *http.Request, tokenCreator *token.JWT) {
				addAuthorization(t, req, tokenCreator, testUserId, false)
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteOneProduct(gomock.Any(), gomock.Eq(product.ID)).
					//Return(product, nil).
					Times(0)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
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

			url := fmt.Sprintf("/api/v1/admin/products/%s", tc.productId)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.auth(t, request, server.TokenCreator())
			server.Router().ServeHTTP(recorder, request)
			tc.response(t, recorder)
		})
	}
}
