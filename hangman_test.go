package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockObject struct {
	mock.Mock
}
type hangmanTestSuite struct {
	word string
	suite.Suite
	mock *MockObject
}

func (suite *hangmanTestSuite) SetupTest() {
	suite.word = "racecar"
	suite.mock = &MockObject{}
}

func (m *MockObject) RenderGame(placeholder []string, chances int, entries map[string]bool) error {
	args := m.Called(placeholder, entries)
	return args.Error(0)
}
func (m *MockObject) getInput() string {
	args := m.Called()
	return args.String(0)
}

func TestHangmanTestSuite(t *testing.T) {
	suite.Run(t, new(hangmanTestSuite))
}

func (suite *hangmanTestSuite) TestPlaySuccess() {
	suite.mock.On("RenderGame", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.mock.On("getInput").Return("racecar")
	result := play(suite.mock, suite.word)
	assert.Equal(suite.T(), result, true)

}

func (suite *hangmanTestSuite) TestPlayFailed() {
	suite.mock.On("RenderGame", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.mock.On("getInput").Return("raceca")
	result := play(suite.mock, suite.word)
	assert.Equal(suite.T(), result, false)

}
