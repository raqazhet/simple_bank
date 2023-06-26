package api

import (
	"context"
	"database/sql"
	"errors"
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

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}
type loginUserResponse struct {
	AccessToken string     `json:"access_token"`
	Login       model.User `json:user`
}

func (s *Server) LoginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := s.store.GetUser(context.Background(), req.Username)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}
	if err := util.CHeckPassword(req.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accsesToken, err := s.tokenMaker.CreateToken(req.Username, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := loginUserResponse{
		AccessToken: accsesToken,
		Login:       user,
	}
	ctx.JSON(http.StatusOK, rsp)
}
