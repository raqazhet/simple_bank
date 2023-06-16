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

func TestGetAccountApi(t *testing.T) {
	account := randomAccount()
	testCases := []struct {
		name          string
		accountid     int
		buildStubs    func(store *mock_storage.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "ok",
			accountid: account.ID,
			buildStubs: func(store *mock_storage.MockStore) {
				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, recorder.Body, account)
			},
		},
	}
	l := jsonlog.Logger{}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mock_storage.NewMockStore(ctrl)
			v.buildStubs(store)

			server := NewServer(store, &l)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/v1/accounts/%d", account.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			v.checkResponse(t, recorder)
		})
	}
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
