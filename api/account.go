package api

import (
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

type getID struct {
	Id int `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccountByID(ctx *gin.Context) {
	var req getID
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.l.PrintError(err, map[string]string{
			"getAccById err": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := s.store.GetAccountById(ctx, req.Id)
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
