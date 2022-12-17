package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockReader struct {
	mock.Mock
}

/**
* Reads a file and returns an bytes array
 */
func (m *MockReader) ReadFile(path string) ([]byte, error) {
	args := m.Called(path)
	return []byte(args.String(0)), args.Error(1)
}
