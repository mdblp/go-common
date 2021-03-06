package shoreline

import (
	"log"
	"strings"
)

type ShorelineMockClient struct {
	ServerToken  string
	Unauthorized bool
	UserID       string
	IsServer     bool
}

func NewMock(token string) *ShorelineMockClient {
	return &ShorelineMockClient{
		ServerToken:  token,
		Unauthorized: false,
		UserID:       "123.456.789",
		IsServer:     true,
	}
}

func (client *ShorelineMockClient) Start() error {
	log.Println("Started mock shoreline client")
	return nil
}

func (client *ShorelineMockClient) Close() {
	log.Println("Close mock shoreline client")
}

func (client *ShorelineMockClient) Login(username, password string) (*UserData, string, error) {
	return &UserData{UserID: client.UserID, Username: username, Emails: []string{username}}, client.ServerToken, nil
}

func (client *ShorelineMockClient) Signup(username, password, email string) (*UserData, error) {
	return &UserData{UserID: client.UserID, Username: username, Emails: []string{email}}, nil
}

func (client *ShorelineMockClient) CheckToken(token string) *TokenData {
	if client.Unauthorized {
		return nil
	}
	return &TokenData{UserID: client.UserID, IsServer: client.IsServer}
}

func (client *ShorelineMockClient) TokenProvide() string {
	return client.ServerToken
}

func (client *ShorelineMockClient) GetUser(userID, token string) (*UserData, error) {
	if userID == "NotFound" {
		return nil, nil
	} else if userID == "WithoutPassword" {
		return &UserData{UserID: userID, Username: "From Mock", Emails: []string{userID}, PasswordExists: false}, nil
	} else if strings.Contains(strings.ToLower(userID), "clinic") {
		return &UserData{UserID: userID, Username: "From Mock", Emails: []string{userID}, PasswordExists: false, Roles: []string{"clinic"}}, nil
	} else {
		return &UserData{UserID: userID, Username: "From Mock", Emails: []string{userID}, PasswordExists: true}, nil
	}
}

func (client *ShorelineMockClient) UpdateUser(userID string, userUpdate UserUpdate, token string) error {
	return nil
}
