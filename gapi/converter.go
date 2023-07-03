package gapi

import (
	db "github.com/aalug/bank-go/db/sqlc"
	"github.com/aalug/bank-go/pb"
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
