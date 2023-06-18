package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"bank/jsonlog"
	"bank/model"
	mock_storage "bank/storage/mock"
	"bank/util"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateUserApi(t *testing.T) {
	user, password := randomUser(t)
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_storage.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"username":        user.Username,
				"hashed_password": password,
				"full_name":       user.Fullname,
				"email":           user.Email,
			},
			buildStubs: func(store *mock_storage.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			l := jsonlog.Logger{}
			ctr1 := gomock.NewController(t)
			defer ctr1.Finish()
			store := mock_storage.NewMockStore(ctr1)
			tc.buildStubs(store)
			server := NewServer(store, &l)
			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			ur1 := "/v1/user/"
			request, err := http.NewRequest(http.MethodPost, ur1, bytes.NewReader(data))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T) (model.User, string) {
	password := util.RandomString(8)
	hashedPassword, err := util.HashedPassword(password)
	require.NoError(t, err)
	arg := model.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		Fullname:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return arg, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user model.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser model.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Fullname, gotUser.Fullname)
	require.Equal(t, user.Email, gotUser.Email)
	require.NotEmpty(t, gotUser.HashedPassword)
}
