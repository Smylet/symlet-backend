package cmd

import (
	"testing"

	"github.com/Smylet/symlet-backend/tests/integration/helpers"
	"github.com/stretchr/testify/suite"
)

type PopulateTestSuite struct {
	suite.Suite
	client *helpers.HttpClient
}

func TestPopulateTestSuite(t *testing.T) {
	suite.Run(t, new(PopulateTestSuite))
}
