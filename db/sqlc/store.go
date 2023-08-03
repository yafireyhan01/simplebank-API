package db

import (
	"context"
	"database/sql"
	"fmt"
)

// store provide all function to execute db queries and transaction
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provide all function to execute db queries and transaction
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore create a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// to execute a generic database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
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

		// adding if smaller id will be executed first to avoid deadlock
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return // same as return account1, account 2, err
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
