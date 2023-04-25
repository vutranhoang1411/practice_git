package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mock_sqlc "github.com/vutranhoang1411/SimpleBank/db/mock"
	db "github.com/vutranhoang1411/SimpleBank/db/sqlc"
	"github.com/vutranhoang1411/SimpleBank/util"
)
func TestGetAccountAPI(t *testing.T) {
	account:=randomAccount()
	tc:=[]struct{
		name string
		accountID int64
		setAuth func(t *testing.T,request *http.Request,server *Server)
		addExpect func (store *mock_sqlc.MockStore)
		checkResponse func (t *testing.T,recorder *httptest.ResponseRecorder)
	}{
		{
			name:"OK",
			accountID: account.ID,
			setAuth: func(t *testing.T, request *http.Request,server *Server) {
				token,_,err:=server.maker.CreateToken(account.Owner,time.Minute)
				require.NoError(t,err)
				request.Header.Add("authorization","bearer "+token)
			},
			addExpect: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.ID)).Times(1).Return(account,nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,recorder.Result().StatusCode,http.StatusOK)
			},
		},
		{
			name:"Not found",
			accountID: 989898,
			setAuth: func(t *testing.T, request *http.Request,server *Server) {
				token,_,err:=server.maker.CreateToken(account.Owner,time.Minute)
				require.NoError(t,err)
				request.Header.Add("authorization","bearer "+token)
			},
			addExpect: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(),gomock.Any()).Times(1).Return(db.Account{},sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusNotFound,recorder.Result().StatusCode)
			},
		},
		{
			name:"Unauthorized",
			accountID: account.ID,
			setAuth: func(t *testing.T, request *http.Request,server *Server) {

			},
			addExpect: func(store *mock_sqlc.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(),gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusUnauthorized,recorder.Result().StatusCode)
			},
		},
		
	}
	//account use for testing
	for i:=range tc{
		t.Run(tc[i].name,func(t *testing.T) {
			//mock db
			ctrl:=gomock.NewController(t)
			defer ctrl.Finish()
			store:=mock_sqlc.NewMockStore(ctrl)

			tc[i].addExpect(store)

			//set up server and request
			server,err:= newTestServer(store)
			require.NoError(t,err);
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/account/%d", tc[i].accountID), nil)

			//set authentication
			tc[i].setAuth(t,request,server)
			server.router.ServeHTTP(recorder, request)
			
			tc[i].checkResponse(t,recorder)
		})
	}


}
func matchingBody(t *testing.T,account db.Account,body *bytes.Buffer){
	var reqBody db.Account
	json.Unmarshal(body.Bytes(),&reqBody)
	require.Equal(t,reqBody,account)
}
func randomAccount() db.Account{
	return db.Account{
		ID:util.RandomNum(1,1000),
		Owner: util.RandomName(),
		Balance: util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}