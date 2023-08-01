package api

import (
	"database/sql"
	"net/http"
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	// same as create account params but without balance (bcs it initial balance should always 0)
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

// create new account
func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// if input data is valid = no error
	// insert a new account into database
	// the input
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// if no error = account successfully created
	ctx.JSON(http.StatusOK, account)
}

// GET ACCOUNT BY ID
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"` // min=1 is for cant get lower than 1
}

// get aspecific account by ID
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		// if account id wasnt found
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// if no error
	ctx.JSON(http.StatusOK, account)
}

// LIST ACCOUNT BY ID
type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`          // min=1 is for cant get lower than 1
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"` // min=1 is for cant get lower than 5
}

// get aspecific account by ID
func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		// if error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// if no error
	ctx.JSON(http.StatusOK, accounts)
}

// DELETE ACCOUNT BY ID
type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"` // min=1 is for cant get lower than 1
}

// delete  account by ID
func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		// if account id wasnt found
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	_ = server.store.DeleteAccount(ctx, req.ID)

	ctx.JSON(http.StatusOK, deleteResponse())

}

// type UpdateAccountRequest struct {
// 	// same as create account params but without balance (bcs it initial balance should always 0)
// 	ID int64 `uri:"id"` // min=1 is for cant get lower than 1
// 	// Owner    string `json:"owner" binding:"required"`
// 	Balance int64 `json:"balance binding:"required,min=1""`
// 	// Currency string `json:"currency" binding:"required,oneof=USD EUR"`
// }

// // update account
// func (server *Server) updateAccount(ctx *gin.Context) {
// 	var req UpdateAccountRequest
// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// 	account, err := server.store.GetAccountForUpdate(ctx, req.ID)
// 	if err != nil {
// 		// if account id wasnt found
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	// if err != nil {
// 	// 	ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 	// }

// 	arg := db.UpdateAccountParams{
// 		Balance: req.Balance,
// 	}

// 	account, err = server.store.UpdateAccount(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	// account2 := db.UpdateAccountParams{
// 	// 	Balance: req.Balance,
// 	// }

// 	// account, _ := server.store.UpdateAccount(ctx, account2)

// 	ctx.JSON(http.StatusOK, account)
// }
