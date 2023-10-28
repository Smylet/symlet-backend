package fixtures

// type UserFixures struct {
// 	baseFixtures
// 	userRepository users.UserRepositoryProvider
// }

// func NewUserFixtures(db db.DBProvider) (*UserFixures, error) {
// 	return &UserFixures{
// 		baseFixtures:   baseFixtures{db: db.GormDB()},
// 		userRepository: users.NewUserRepository(db.GormDB()),
// 	}, nil
// }

// func (f UserFixures) CreateUser(ctx context.Context, arg users.CreateUserTxParams) (users.User, error) {
// 	user, err := f.userRepository.CreateUserTx(ctx, arg)
// 	if err != nil {
// 		log.Fatal().Err(err).Str("email", arg.Email).Msg("failed to create user")
// 		return users.User{}, err
// 	}
// 	return user.User, nil
// }

// func (f UserFixures) FindUser(ctx context.Context, arg users.FindUserParams) (users.User, error) {
// 	user, err := f.userRepository.FindUser(ctx, arg)
// 	if err != nil {
// 		log.Fatal().Err(err).Str("email", arg.Email).Msg("failed to find user")
// 		return users.User{}, err
// 	}
// 	return user.User, nil
// }
