package api

import (
	"database/sql"
	"errors"
	"net/http"

	"bank/model"

	"github.com/gin-gonic/gin"
)

func (s *Server) CreateAccount(ctx *gin.Context) {
	var input model.CreateAccountParams
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := model.CreateAccountParams{
		Owner:    input.Owner,
		Currency: input.Currency,
		Balance:  0,
	}
	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	Id int `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccountByID(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		s.l.PrintError(err, map[string]string{
			"id": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// accountId := strings.TrimPrefix(ctx.Request.URL.Path, "/v1/accounts/")
	// id, err := strconv.Atoi(accountId)
	// if err != nil {
	// 	s.l.PrintError(err, map[string]string{
	// 		"id": err.Error(),
	// 	})
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }
	account, err := s.store.GetAccountById(ctx, req.Id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// type listAccountRequest struct {
// 	PageId   int `form:"page_id" binding:"required,min=1"`
// 	PageSize int `form:"page_id" binding:"required,min=5,max=10"`
// }

func (s *Server) getAllAccounts(ctx *gin.Context) {
	accounts, err := s.store.GetAllAccounts(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
