package gapi

import (
	"bank/model"
	"bank/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user model.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.Fullname,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
