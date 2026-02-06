package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ClientSuite struct {
	suite.Suite
	base *BaseTestSuite
}

func (s *ClientSuite) TestHappyPath() {
	s.T().Run("register flow", s.RegisterFlow)
}

func (s *ClientSuite) RegisterFlow(t *testing.T) {
	ok, err := s.base.client.Health()
	s.Require().NoError(err)
	s.Require().True(ok)
}
