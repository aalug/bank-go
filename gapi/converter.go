package gapi

import (
	db "github.com/aalug/go-bank/db/sqlc"
	"github.com/aalug/go-bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// convertUser converts a db.User object to a UserResponse object
func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
