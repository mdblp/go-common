package opa

import (
	"fmt"
	"net/http"
)

type opaAuthCall struct {
	auth *Authorization
	err  error
}

// MockClient The mocked interface to portal-api.
//
type MockClient struct {
	nextOpaAuthCall map[string]*opaAuthCall
}

// NewMock create a new portal mock client
func NewMock() *MockClient {
	return &MockClient{
		nextOpaAuthCall: make(map[string]*opaAuthCall),
	}
}

// SetMockOpaAuth Set the result for the next MockPatientConfig calls
//
// - token: The token string for which the response will be
//
// - pc: The PatientConfig to return or nil
//
// - err: The error to return or nil
func (client *MockClient) SetMockOpaAuth(key string, auth *Authorization, err error) {
	client.nextOpaAuthCall[key] = &opaAuthCall{
		auth: auth,
		err:  err,
	}
}

//GetMockedAuth returns Authorization struct from external packages (opaResult stuct is "private")
func (client *MockClient) GetMockedAuth(authorized bool, data map[string]interface{}, route string) Authorization {
	return Authorization{
		Result: &opaResult{
			Authorized: authorized,
			Data:       data,
			Route:      route,
		},
	}
}

// GetOpaAuth mock the GetPatientConfig call
func (client *MockClient) GetOpaAuth(req *http.Request, data map[string]interface{}) (*Authorization, error) {
	key := req.Host + req.URL.Path
	pcc, ok := client.nextOpaAuthCall[key]
	if !ok {
		return nil, fmt.Errorf("Unknown response code[404] from service[http://opa/%s]", routeAuth)
	}
	return pcc.auth, pcc.err
}
