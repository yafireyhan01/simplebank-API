package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestTransferTx(t *testing.T) {
// 	store := NewStore(testDB)

// 	account1 := createRandomAccount(t)
// 	account2 := createRandomAccount(t)

// 	fmt.Println(">>before: ", account1.Balance, account2.Balance)

// 	// run n concurrent transfer transaction
// 	n := 5
// 	amount := int64(10)

// 	// use make() keyword to create the channel
// 	// send them back to the main go routine that our test is running on
// 	// channel is designed to connect concurrent Go routines
// 	errs := make(chan error)
// 	// in our case we need 1 channel to recieve the errors,
// 	// and 1 other channel to recieve the TransferTxResult
// 	results := make(chan TransferTxResult)
// 	// ======================

// 	for i := 0; i < n; i++ {
// 		// print out some logs to see wich transaction is calling wich query and in wich order
// 		// txName := fmt.Sprintf("tx %d", i+1) // for debug purpose

// 		// go func to start new goroutine
// 		go func() {
// 			ctx := context.Background()
// 			result, err := store.TransferTx(ctx, TransferTxParams{
// 				FromAccountID: account1.ID,
// 				ToAccountID:   account2.ID,
// 				Amount:        amount,
// 			})

// 			// send error to the errors channel using this arrow operation
// 			// channel should be on the left, and data sent should be on the right of the arrow operator
// 			errs <- err
// 			results <- result
// 		}()
// 	}

// 	// check errors and results from outside
// 	// CHECK results
// 	existed := make(map[int]bool)

// 	for i := 0; i < n; i++ {
// 		err := <-errs
// 		require.NoError(t, err)

// 		result := <-results
// 		require.NotEmpty(t, result)

// 		// as result contains several object inside, lets verify each of them
// 		// CHECK transfer
// 		transfer := result.Transfer
// 		require.NotEmpty(t, transfer)
// 		// the fromaccountID field of transfer should equal to account1.ID
// 		require.Equal(t, account1.ID, transfer.FromAccountID)
// 		require.Equal(t, account2.ID, transfer.ToAccountID)
// 		require.Equal(t, amount, transfer.Amount) // transfer.amount = input amount
// 		require.NotZero(t, transfer.ID)           // The ID transfer shoudnt be zero cs its an auto increment field
// 		require.NotZero(t, transfer.CreatedAt)    // The createdAt shoudnt be zero cs database expected to fill in the default value

// 		// check if the transfer record is really created on the database
// 		_, err = store.GetTransfer(context.Background(), transfer.ID)
// 		// if the transfer really exist, the function shouldnt return an error
// 		require.NoError(t, err)

// 		// CHECK entries
// 		fromEntry := result.FromEntry
// 		require.NotEmpty(t, fromEntry)
// 		require.Equal(t, account1.ID, fromEntry.AccountID)
// 		require.Equal(t, -amount, fromEntry.Amount)
// 		require.NotZero(t, fromEntry.ID)
// 		require.NotZero(t, fromEntry.CreatedAt)

// 		// try to get account entry from database, to make sure that its really got created
// 		_, err = store.GetEntry(context.Background(), fromEntry.ID)
// 		require.NoError(t, err)

// 		// =============
// 		toEntry := result.ToEntry
// 		require.NotEmpty(t, toEntry)
// 		require.Equal(t, account2.ID, toEntry.AccountID)
// 		require.Equal(t, amount, toEntry.Amount)
// 		require.NotZero(t, toEntry.ID)
// 		require.NotZero(t, toEntry.CreatedAt)

// 		// try to get account entry to database, to make sure that its really got created
// 		_, err = store.GetEntry(context.Background(), toEntry.ID)
// 		require.NoError(t, err)

// 		// CHECK ACCOUNT
// 		fromAccount := result.FromAccount
// 		require.NotEmpty(t, fromAccount)
// 		require.Equal(t, account1.ID, fromAccount.ID)
// 		// where money is going into
// 		toAccount := result.ToAccount
// 		require.NotEmpty(t, toAccount)
// 		require.Equal(t, account2.ID, toAccount.ID)

// 		// CHECK ACCOUNT BALANCE
// 		fmt.Println(">>tx: ", fromAccount.Balance, toAccount.Balance)
// 		diff1 := account1.Balance - fromAccount.Balance
// 		diff2 := toAccount.Balance - account2.Balance
// 		require.Equal(t, diff1, diff2)
// 		require.True(t, diff1 > 0)         // should be positive
// 		require.True(t, diff1%amount == 0) // 1*amount, 2*amount, 3*amount, ...., n*amount

// 		k := int(diff1 / amount)
// 		require.True(t, k >= 1 && k <= n) // k antara 1 dan n (positif)
// 		require.NotContains(t, existed, k)
// 		existed[k] = true // set existed k to true

// 		// CHECK THE FINAL UPDATED BALANCE
// 		updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
// 		require.NoError(t, err)
// 		updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
// 		require.NoError(t, err)

// 		fmt.Println(">>after: ", updatedAccount1.Balance, updatedAccount2.Balance)

// 		require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
// 		require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
// 	}
// }

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
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

		// check balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
