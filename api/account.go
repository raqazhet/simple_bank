package api

import (
	"net/http"
	"strconv"
	"strings"

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

// type getAccountRequest struct {
// 	ID int `uri:"id"`
// }

func (s *Server) getAccountByID(ctx *gin.Context) {
	accountId := strings.TrimPrefix(ctx.Request.URL.Path, "/v1/accounts/")
	id, err := strconv.Atoi(accountId)
	if err != nil {
		return
	}
	account, err := s.store.GetAccountById(ctx, id)
	if err != nil {
		s.l.PrintError(err, map[string]string{
			"account err": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, account)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
