package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/vutranhoang1411/SimpleBank/util"
)

func createRandomUser(t *testing.T) User{
	arg:=CreateUserParams{
		ID:uuid.New().String()[2:18],
		Name:util.RandomName(),
		Password: "secret",
		Email: util.RandomEmail(),
	}
	user,err:=testQueries.CreateUser(context.Background(),arg)
	require.NoError(t,err);
	require.NotEmpty(t,user)

	require.Equal(t,arg.ID,user.ID)
	require.Equal(t,arg.Name,user.Name)
	require.Equal(t,arg.Email,user.Email)
	require.Equal(t,arg.Password,user.Password)
	return user

}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
func TestGetUser(t *testing.T) {
	user1:=createRandomUser(t);
	user2,err:=testQueries.GetUser(context.Background(),user1.ID);
	require.NoError(t,err);
	require.Equal(t,user1.ID,user2.ID)

}