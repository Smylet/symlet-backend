//go:build integration

package users

import (
	"context"
	"testing"

	"github.com/rs/zerolog/log"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/tests/integration/fixtures"
	"github.com/Smylet/symlet-backend/tests/integration/helpers"
	"github.com/Smylet/symlet-backend/utilities/db"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CreateUserTestSuite struct {
	suite.Suite
	client       *helpers.HttpClient
	userFixtures *fixtures.UserFixures
}

func TestCreateUserTestSuite(t *testing.T) {
	suite.Run(t, new(CreateUserTestSuite))
}

func (s *CreateUserTestSuite) SetupSuite() {
	config, err := utils.LoadConfig()
	if err != nil {
		s.T().Fatal("Error loading config: ", err)
		return // No need to continue
	}

	db, err := db.GetDB(config)
	if err != nil {
		log.Fatal().Msg("Error establishing database connection: ")
		return // No need to continue
	}
	serviceUrl := config.HTTPServerAddress
	s.client = helpers.NewUserApiClient(serviceUrl)
	fixtures, err := fixtures.NewUserFixtures(db)
	if err != nil {
		log.Fatal().Msg("Error loading user fixtures: ")
		return // No need to continue
	}
	s.userFixtures = fixtures
}

func (s *CreateUserTestSuite) TestOk() {
	defer func() { assert.Nil(s.T(), s.userFixtures.UnloadFixtures()) }()

	ctx := context.Background()

	tests := []struct {
		name     string
		request  users.CreateUserReq
		expected utils.SuccessMessage
	}{
		// ... (existing test cases) ...
	}

	for _, test := range tests {
		test := test // Capture range variable.
		s.Run(test.name, func() {
			var response utils.SuccessMessage
			err := s.client.DoPostRequest("/register", test.request, &response)
			assert.Nil(s.T(), err)
			assert.Equal(s.T(), test.expected.Msg, response.Msg)

			// Check if the Data is of type map[string]interface{}
			dataMap, ok := response.Data.(map[string]interface{})
			assert.True(s.T(), ok, "Response Data is not of type map[string]interface{}")

			if ok { // Only proceed if type assertion was successful
				for key, expectedValue := range test.expected.Data.(map[string]interface{}) {
					actualValue, exists := dataMap[key]
					assert.True(s.T(), exists, "Key %s does not exist in response data", key)
					assert.Equal(s.T(), expectedValue, actualValue, "Value mismatch for key %s", key)
				}
			}

			// Check user in the database using the FindUser method
			user, err := s.userFixtures.FindUser(ctx, users.FindUserParams{
				User: users.User{
					Username: test.request.Username,
				},
			})

			assert.Nil(s.T(), err, "User not found in the database")
			assert.Equal(s.T(), test.request.Username, user.Username, "Username in DB does not match the request")
			assert.Equal(s.T(), test.request.Email, user.Email, "Email in DB does not match the request")
		})
	}
}
