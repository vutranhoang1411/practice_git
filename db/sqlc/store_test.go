package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)
func TestTransferTx(t *testing.T) {
	store:=NewStore(testDB)

	account1:=createRandomAccount(t)
	account2:=createRandomAccount(t)
	n:=2
	amount:=int64(10)
	errs:=make(chan error)
	results:=make(chan TransferTxResult)
	for i:=0;i<n;i++{
		go func(){
			result,err:=store.TransferTx(context.Background(),TransferTxParams{
				FromAccountID:	account1.ID,
				ToAccountID: 	account2.ID,
				Amount: 		amount,
			})
			errs<-err
			results<-result
		}()
	}
	
	for i:=0;i<n;i++{
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		diff1:=account1.Balance-fromAccount.Balance;
		diff2:=toAccount.Balance-account2.Balance;
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
	}
}
func TestTransferTXDeadlock(t *testing.T) {
	store:=NewStore(testDB);
	account1:=createRandomAccount(t);
	account2:=createRandomAccount(t);

	n:=6;
	amount:=int64(10);

	errs:=make(chan error)
	for i:=0;i<n;i++{
		var fromAccount,toAccount Account;
		fromAccount=account1;
		toAccount=account2;
		if (i%2==1){
			fromAccount=account2
			toAccount=account1
		}
		go func (){
			_,err:=store.TransferTx(context.Background(),TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID: toAccount.ID,
				Amount: amount,
			})
			errs<-err
		}()
	}

	for i:=0;i<n;i++{
		err:=<-errs
		require.NoError(t,err)
	}
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
func TestWithDrawTx(t *testing.T) {
	store:=NewStore(testDB);
	n:=2;
	amount:=int64(10);
	account1:=createRandomAccount(t);
	account2:=createRandomAccount(t);
	fmt.Print("before: ",account1.Balance,account2.Balance)
	errs:=make(chan error);
	for i:=0;i<n;i++{
		go func(i int) {
			var err error;
			if (i%2==0){
				_,err=store.TransferTx(context.Background(),TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID: account2.ID,
					Amount: amount,
				})
				errs<-err;
			}else{
				_,err:=store.WithDrawTx(context.Background(),WithDrawMoneyArg{
					AccountID: account2.ID,
					Amount: amount,
				})
				errs<-err
			}
		}(i)
	}
	for i:=0;i<n;i++{
		err:=<-errs;
		require.NoError(t,err)
	}
	fromAccount,err:=store.GetAccount(context.Background(),account1.ID);
	require.NoError(t,err)
	require.NotEmpty(t,fromAccount)
	require.Equal(t,fromAccount.Balance,account1.Balance-int64(n)/2*amount);


	toAccount,err:=store.GetAccount(context.Background(),account2.ID)
	require.NoError(t,err)
	require.NotEmpty(t,toAccount);
	require.Equal(t,account2.Balance,toAccount.Balance);
}