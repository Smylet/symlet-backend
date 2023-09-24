package users

import (
	"testing"

	"github.com/rs/zerolog/log"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/tests/integration/fixtures"
	"github.com/Smylet/symlet-backend/tests/integration/helpers"
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
		log.Fatal().Err(err).Msg("Error loading config")
	}
	serviceUrl := config.HTTPServerAddress
	s.client = helpers.NewUserApiClient(serviceUrl)
	fixtures, err := fixtures.NewUserFixtures()
	s.Require().NoError(err)
	s.userFixtures = fixtures

}

func (s *CreateUserTestSuite) TestOk() {
	defer func() { assert.Nil(s.T(), s.userFixtures.UnloadFixtures()) }()

	tests := []struct {
		name     string
		request  users.CreateUserReq
		expected utils.SuccessMessage
	}{
		{
			name: "CreateValidUser",
			request: users.CreateUserReq{
				Username: "test",
				Email:    "test@gmail.com",
				Password: "PA1SSWORd1",
			},
			expected: utils.SuccessMessage{
				Msg: "User created successfully",
				Data: map[string]interface{}{
					"username": "test",
					"email":    "test@gmail.com",
				},
			},
		},
	}
	for _, test := range tests {
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
		})
	}

}
