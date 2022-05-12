package auth

import (
	"net/http"

	"github.com/mdblp/shoreline/token"
)

type ClientMock struct {
	ServerToken  string
	Unauthorized bool
	UserID       string
	IsServer     bool
}

func NewMock(token string) *ClientMock {
	return &ClientMock{
		ServerToken:  token,
		Unauthorized: false,
		UserID:       "123.456.789",
		IsServer:     true,
	}
}

func (client *ClientMock) Authenticate(req *http.Request) *token.TokenData {
	if client.Unauthorized {
		return nil
	}

	if sessionToken := req.Header.Get("x-tidepool-session-token"); sessionToken != "" {
		return &token.TokenData{UserId: client.UserID, IsServer: client.IsServer}
	} else if req.Header.Get("authorization") != "" {
		return &token.TokenData{UserId: client.UserID, IsServer: client.IsServer}
	}
	return nil
}
