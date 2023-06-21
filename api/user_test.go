package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"bank/model"
	mock_storage "bank/storage/mock"
	"bank/util"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func EqCreateUserParams(arg model.User, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

type Matcher interface {
	// Matches returns whether x is a match
	Matches(x interface{}) bool
	// String describes what the matcher matches
	String() string
}
type eqCreateUserParamsMatcher struct {
	arg      model.User
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(model.User)
	if !ok {
		return false
	}
	err := util.CHeckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func TestCreateUserApi(t *testing.T) {
	user, password := randomUser(t)
	// hashedPassword, err := util.HashedPassword(password)
	// require.NoError(t, err)
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
				arg := model.User{
					Username: user.Username,
					// HashedPassword: hashedPassword,
					Fullname: user.Fullname,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
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
			// l := jsonlog.Logger{}
			ctr1 := gomock.NewController(t)
			defer ctr1.Finish()
			store := mock_storage.NewMockStore(ctr1)
			tc.buildStubs(store)
			server := newTestServer(t, store)
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
