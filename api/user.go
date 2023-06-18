package api

import (
	"context"
	"log"
	"net/http"

	"bank/model"
	"bank/util"

	"github.com/gin-gonic/gin"
)

func (s *Server) CreateUser(ctx *gin.Context) {
	var req model.User
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("createUser bindJson err: %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := util.HashedPassword(req.HashedPassword)
	if err != nil {
		s.l.PrintError(err, map[string]string{
			"HashPassowrd": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	req.HashedPassword = hashedPassword
	user, err := s.store.CreateUser(context.Background(), req)
	if err != nil {
		s.l.PrintError(err, map[string]string{
			"userPsql error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}
