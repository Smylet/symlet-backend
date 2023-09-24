package fixtures

import (
	"context"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/db"
	"github.com/Smylet/symlet-backend/utilities/utils"
)

type UserFixures struct {
	baseFixtures
	userRepository users.UserRepositoryProvider
}

func NewUserFixtures() (*UserFixures, error) {
	config, err := utils.LoadConfig()
	if err != nil {
		return nil, err
	}

	db, err := db.GetDB(config)
	if err != nil {
		return nil, err
	}
	return &UserFixures{
		baseFixtures:   baseFixtures{db: db.GormDB()},
		userRepository: users.NewUserRepository(db.GormDB()),
	}, nil
}

func (f UserFixures) CreateUser(ctx context.Context, arg users.CreateUserTxParams) (users.User, error) {
	user, err := f.userRepository.CreateUserTx(ctx, arg)
	if err != nil {
		return users.User{}, err
	}
	return user.User, nil
}
