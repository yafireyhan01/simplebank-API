package db

import (
	"context"
	"database/sql"
	"fmt"
)

// store provide all function to execute db queries and transaction
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore create a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// to execute a generic database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		// rbErr = Rollback Error
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		// if rollback is successful
		return err // original transaction error
	}
	// if all operations in the transaction are successful
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`     // the created transfer record
	FromAccount Account  `json:"from_account"` // after the balance is updated
	ToAccount   Account  `json:"to_account"`   // after the balance is updated
	FromEntry   Entry    `json:"from_entry"`   // the entry of the from account wich records that money is moving out
	ToEntry     Entry    `json:"to_entry"`     // the entry of the from account wich records that money is moving in
}

// {}{} means we are creating a new empty object of that type
// var txKey = struct{}{} // for debug purpose

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update accounts balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// func(q *Queries) error {} = callback function
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// txName := ctx.Value(txKey) // for debug purpose

		// fmt.Println(txName, "create transfer") // for debug purpose
		// transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
			// the output transfer will be saved to result.transfer
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "create entry 1") // for debug purpose
		// account entries
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // because money is moving out from this account
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "create entry 2")// for debug purpose
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount, // because money is moving into this account
		})
		if err != nil {
			return err
		}

		// update accounts balance
		// get account from DB => add or substract amount of money from its balance => update it back to DB
		// 1. get account
		// fmt.Println(txName, "get account1 for update") // for debug purpose
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }

		// fmt.Println(txName, "update account2 balance") // for debug purpose
		//else
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount, // - bcs the money were moving out
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "get account2 for update") // for debug purpose
		// move money to account2 (account2 balance increasing)
		// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return err
		// }

		// fmt.Println(txName, "update account2 balance") // for debug purpose
		//else
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
