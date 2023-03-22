package auth

import (
	"net/http"

	"github.com/mdblp/shoreline/token"
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func NewMock() *ClientMock {
	return &ClientMock{}
}

func (client *ClientMock) Authenticate(req *http.Request) *token.TokenData {
	args := client.Called(req)
	if args.Get(0) == nil {
		return nil
	} else {
		return args.Get(0).(*token.TokenData)
	}
}
