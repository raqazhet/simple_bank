package gapi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"bank/model"
	"bank/pb"
	"bank/util"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (srv *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashedPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
	}
	arg := model.User{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Fullname:       req.FullName,
		Email:          req.Email,
	}
	user, err := srv.store.CreateUser(context.Background(), arg)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			switch pqerr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}
	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

func (srv *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := srv.store.GetUser(ctx, req.Username)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, status.Errorf(codes.NotFound, "user not fount %s", err)
		default:
			return nil, status.Errorf(codes.Internal, "failed to find user %s", err)
		}
	}
	err = util.CHeckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password %s", err)
	}
	accessToken, accsesPayload, err := srv.tokenMaker.CreateToken(user.Username, srv.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token %s", err)
	}
	meta := srv.extractMetadata(ctx)
	fmt.Println(meta)
	rsp := &pb.LoginUserResponse{
		User:                 convertUser(user),
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accsesPayload.ExpiredAt),
	}
	return rsp, nil
}
