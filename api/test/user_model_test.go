package test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/go-faker/faker/v4"
)

func TestUser(t *testing.T) {
	log.Println(os.Getwd())



	userReq := users.CreateUserReq{
		Username: "test",
		Email: faker.Email(),
		Password: "test",
	}
	// Create user
	ctx := context.Background()
	_, err := users.CreateUserTx(ctx, DB, users.CreateUserTxParams{
		CreateUserReq: userReq,
		AfterCreate: func(user users.User) error {
			return nil
		},
	})
	
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	// Check if user was created
	var user users.User
	DB.First(&user, "username = ?", userReq.Username)
	if user.Username != userReq.Username {
		t.Errorf("Expected username to be %v, got %v", userReq.Username, user.Username)
	}
	if user.Email != userReq.Email {
		t.Errorf("Expected email to be %v, got %v", userReq.Email, user.Email)
	}

}

func TestMain(m *testing.M) {
    exitCode := RunTests(m, "../../env")
    os.Exit(exitCode)
}