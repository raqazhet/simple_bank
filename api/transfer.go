package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"bank/storage"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountId int    `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int    `json:"to_account_id" binding:"required,min=1"`
	Amount        int    `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !s.validAccount(ctx, req.FromAccountId, req.Currency) {
		return
	}
	if !s.validAccount(ctx, req.ToAccountId, req.Currency) {
		return
	}
	arg := storage.TransferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID:   req.ToAccountId,
		Amount:        req.Amount,
	}
	res, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (srv *Server) validAccount(ctx *gin.Context, accountId int, currency string) bool {
	account, err := srv.store.GetAccountById(ctx, accountId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
