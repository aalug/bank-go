package db

import (
	"context"
	"database/sql"
	"github.com/aalug/bank-go/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)

	params := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, params.Username, user.Username)
	require.Equal(t, params.HashedPassword, user.HashedPassword)
	require.Equal(t, params.FullName, user.FullName)
	require.Equal(t, params.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

// TestCreateUser tests the create user function
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

// TestGetUser tests the get user function
func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)
	params := UpdateUserParams{
		Username: user1.Username,
		HashedPassword: sql.NullString{
			String: utils.RandomString(6),
			Valid:  true,
		},
		FullName: sql.NullString{
			String: utils.RandomOwner(),
			Valid:  true,
		},
		Email: sql.NullString{
			String: utils.RandomEmail(),
			Valid:  true,
		},
	}

	user2, err := testQueries.UpdateUser(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, params.HashedPassword.String, user2.HashedPassword)
	require.Equal(t, params.FullName.String, user2.FullName)
	require.Equal(t, params.Email.String, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
