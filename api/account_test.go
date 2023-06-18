package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"bank/jsonlog"
	"bank/model"
	mock_storage "bank/storage/mock"
	"bank/util"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAccountAPI(t *testing.T) {
	l := jsonlog.Logger{}
	account := randomAccount()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mock_storage.NewMockStore(ctrl)
	store.EXPECT().
		GetAccountById(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)
	server := NewServer(store, &l)
	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/v1/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)
}

func TestCreateAccountApi(t *testing.T) {
	l := jsonlog.Logger{}
	account := randomAccount()
	account.Owner = "Razaq"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mock_storage.NewMockStore(ctrl)
	store.EXPECT().
		CreateAccount(gomock.Any(), gomock.Any()).
		Times(1).
		Return(account, nil)
	server := NewServer(store, &l)
	recorder := httptest.NewRecorder()
	data, err := json.Marshal(account)
	require.NoError(t, err)
	url := "/v1/accounts/"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)
}

func TestGetAllAccountApi(t *testing.T) {
	l := jsonlog.Logger{}
	accounts := []model.Account{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mock_storage.NewMockStore(ctrl)
	store.EXPECT().
		GetAllAccounts(gomock.Any()).Times(1).
		Return(accounts, nil)
	server := NewServer(store, &l)
	recorder := httptest.NewRecorder()
	ur1 := "/v1/accounts"
	request, err := http.NewRequest(http.MethodGet, ur1, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.NotEmpty(t, accounts)
}

func randomAccount() model.Account {
	return model.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account model.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var goAccount model.Account
	err = json.Unmarshal(data, &goAccount)
	require.NoError(t, err)
	require.Equal(t, account, goAccount)
}
