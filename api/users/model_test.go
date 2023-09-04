package users

import (
	"context"
	"os"
	"testing"

	"github.com/Smylet/symlet-backend/api/test"
	"github.com/go-faker/faker/v4"
)

func TestUser(t *testing.T) {
	test.DB.AutoMigrate(
		User{},
		Profile{},
	)

	userReq := CreateUserReq{
		Username: "test",
		Email: faker.Email(),
		Password: "test",
	}
	// Create user
	ctx := context.Background()
	_, err := CreateUserTx(ctx, test.DB, CreateUserTxParams{
		CreateUserReq: userReq,
		AfterCreate: func(user User) error {
			return nil
		},
	})
	
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	// Check if user was created
	var user User
	test.DB.First(&user, "username = ?", userReq.Username)
	if user.Username != userReq.Username {
		t.Errorf("Expected username to be %v, got %v", userReq.Username, user.Username)
	}
	if user.Email != userReq.Email {
		t.Errorf("Expected email to be %v, got %v", userReq.Email, user.Email)
	}

}

func TestMain(m *testing.M) {
    exitCode := test.RunTests(m)
    os.Exit(exitCode)
}