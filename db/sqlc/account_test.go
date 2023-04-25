package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vutranhoang1411/SimpleBank/util"
)
func createRandomAccount(t *testing.T)Account{
	user:=createRandomUser(t);
	arg:=CreateAccountParams{
		Owner: user.ID,
		Balance: util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
	account,err:=testQueries.CreateAccount(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,account)
	require.Equal(t,account.Owner,arg.Owner)
	require.Equal(t,account.Balance,arg.Balance)
	require.Equal(t,account.Currency,arg.Currency)

	require.NotEmpty(t,account.ID)
	require.NotEmpty(t,account.CreatedAt)
	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t);
}
func TestGetAccount(t *testing.T) {
	account1:=createRandomAccount(t);
	account2,err:=testQueries.GetAccount(context.Background(),account1.ID);

	require.NoError(t,err)
	require.NotEmpty(t,account2)

	//equal
	require.Equal(t,account1.ID,account2.ID)
	require.Equal(t,account1.Balance,account2.Balance);
	require.Equal(t,account1.Currency,account2.Currency);
	require.Equal(t,account1.Owner,account2.Owner);
	require.WithinDuration(t,account1.CreatedAt,account2.CreatedAt,time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1:=createRandomAccount(t);
	arg:=UpdateAccountParams{
		ID:account1.ID,
		Balance: util.RandomBalance(),
	}
	account2,err:=testQueries.UpdateAccount(context.Background(),arg);

	require.NoError(t,err);
	require.NotEmpty(t,account2);

	require.Equal(t,account1.ID,account2.ID);
	require.Equal(t,account1.Owner,account2.Owner);
	require.Equal(t,account1.Currency,account2.Currency);
	require.WithinDuration(t,account1.CreatedAt,account2.CreatedAt,time.Second);
	require.Equal(t,account2.Balance,arg.Balance);
}

func TestDeleteAccount(t *testing.T) {
	account:=createRandomAccount(t);
	err:=testQueries.DeleteAccount(context.Background(),account.ID);
	require.NoError(t,err);

	account2,err:=testQueries.GetAccount(context.Background(),account.ID);
	require.Error(t,err)
	require.Empty(t,account2)
}

func TestAddAccountBalance(t *testing.T) {
	account1:=createRandomAccount(t);
	arg:=AddAccountBalanceParams{
		Amount:util.RandomBalance(),
		ID:account1.ID,
	}
	account2,err:=testQueries.AddAccountBalance(context.Background(),arg)

	require.NotEmpty(t,account2);
	require.NoError(t,err)

	require.Equal(t,account1.ID,account2.ID)
	require.Equal(t,account1.Balance+arg.Amount,account2.Balance);
	require.Equal(t,account1.Currency,account2.Currency);
	require.Equal(t,account1.Owner,account2.Owner);
	require.WithinDuration(t,account1.CreatedAt,account2.CreatedAt,time.Second)
}

func TestListAccounts(t *testing.T) {
	time:=5;
	owner:=createRandomUser(t);
	for i:=0;i<time;i++{
		testQueries.CreateAccount(context.Background(),CreateAccountParams{
			Owner: owner.ID,
			Balance: util.RandomBalance(),
			Currency: util.RandomCurrency(),
		})
	}
	accounts,err:=testQueries.ListAccounts(context.Background(),ListAccountsParams{
		Owner: owner.ID,
		Limit: 2,
		Offset: 3,
	})

	require.NoError(t,err);
	require.NotEmpty(t,accounts)
	require.Equal(t,len(accounts),2)

	for i:=0;i<2;i++{
		require.NotEmpty(t,accounts[i])
		require.Equal(t,accounts[i].Owner,owner.ID)
	}
}