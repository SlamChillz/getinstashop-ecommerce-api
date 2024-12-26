package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mockdb "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/mock"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	testCases := []struct {
		name     string
		body     gin.H
		stubs    func(store *mockdb.MockStore)
		response func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Valid Credentials",
			body: gin.H{
				"email":    "test@gmail.com",
				"password": "password123",
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "No Email",
			body: gin.H{
				"password": "password123",
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "No Password",
			body: gin.H{
				"email": "test@gmail.com",
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "No Email, No Password",
			body: gin.H{},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := "/api/v1/auth/register"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
			require.NoError(t, err)
			server.Router().ServeHTTP(recorder, req)
			tc.response(t, recorder)
		})
	}
}

func TestLogin(t *testing.T) {
	pass, hashPass := randomPassword(t)
	user := db.GetUserByIdRow{
		ID:       uuid.UUID{},
		Email:    "test@gmail.com",
		Password: hashPass,
		Admin:    false,
	}
	testCases := []struct {
		name     string
		body     gin.H
		stubs    func(store *mockdb.MockStore)
		response func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Valid Login Details",
			body: gin.H{
				"email":    "test@gmail.com",
				"password": pass,
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserById(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).Times(1)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"email":    "test@gmail.com",
				"password": pass,
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserById(gomock.Any(), gomock.Eq(user.Email)).Return(db.GetUserByIdRow{}, sql.ErrConnDone).Times(1)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Incorrect Password",
			body: gin.H{
				"email":    "test@gmail.com",
				"password": "testuserpassword",
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserById(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).Times(1)
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := "/api/v1/auth/login"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
			require.NoError(t, err)
			server.Router().ServeHTTP(recorder, req)
			tc.response(t, recorder)
		})
	}
}

func randomPassword(t *testing.T) (string, string) {
	pass := "password123"
	hashPass, err := utils.HashPassword(pass)
	require.NoError(t, err)
	return pass, hashPass
}
