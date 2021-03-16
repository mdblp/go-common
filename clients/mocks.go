package clients

import (
	"encoding/json"
	"fmt"
)

type (
	SeagullMock struct {
		nextCollectionCall map[string]*SeagullCollectionCalls
	}
	SeagullCollectionCalls struct {
		result string
		err    error
	}
)

//A mock of the Gatekeeper interface
func NewGatekeeperMock(expectedPermissions Permissions, expectedError error) *GatekeeperMock {
	return &GatekeeperMock{expectedPermissions, expectedError, []string{}}
}
func (mock *GatekeeperMock) SetExpected(expectedPermissions Permissions, expectedError error) {
	mock.expectedPermissions = expectedPermissions
	mock.expectedError = expectedError
}

func (mock *GatekeeperMock) UserInGroup(userID, groupID string) (Permissions, error) {
	if mock.expectedPermissions != nil || mock.expectedError != nil {
		return mock.expectedPermissions, mock.expectedError
	}
	return Permissions{"root": Allowed}, nil
}

func (mock *GatekeeperMock) UsersInGroup(groupID string) (UsersPermissions, error) {
	if mock.expectedPermissions != nil || mock.expectedError != nil {
		return UsersPermissions{groupID: mock.expectedPermissions}, mock.expectedError
	}
	return UsersPermissions{groupID: Permissions{groupID: Allowed}}, nil
}

func (mock *GatekeeperMock) GroupsForUser(userID string) (UsersPermissions, error) {
	if mock.expectedPermissions != nil || mock.expectedError != nil {
		if len(mock.UserIDs) > 0 {
			perms := make(map[string]Permissions)
			for _, user := range mock.UserIDs {
				perms[user] = mock.expectedPermissions
			}
			perms[userID] = Permissions{"root": Allowed}
			return perms, mock.expectedError
		}
		return UsersPermissions{userID: mock.expectedPermissions}, mock.expectedError
	}
	return UsersPermissions{userID: Permissions{userID: Allowed}}, nil
}

func (mock *GatekeeperMock) SetPermissions(userID, groupID string, permissions Permissions) (Permissions, error) {
	return Permissions{"root": Allowed}, nil
}

//A mock of the Seagull interface
func NewSeagullMock() *SeagullMock {
	return &SeagullMock{
		nextCollectionCall: make(map[string]*SeagullCollectionCalls),
	}
}
func (mock *SeagullMock) SetMockNextCollectionCall(key string, expectedResult string, expectedError error) {
	mock.nextCollectionCall[key] = &SeagullCollectionCalls{
		result: expectedResult,
		err:    expectedError,
	}
}
func (mock *SeagullMock) GetPrivatePair(userID, hashName, token string) *PrivatePair {
	return &PrivatePair{ID: "mock.id.123", Value: "mock value"}
}

func (mock *SeagullMock) GetCollection(userID, collectionName, token string, v interface{}) error {
	response, ok := mock.nextCollectionCall[userID+collectionName]
	if !ok {
		return fmt.Errorf("Unknown response code[404] from seagull.getCollection")
	}
	if response.err != nil {
		return response.err
	}

	json.Unmarshal([]byte(response.result), &v)
	return nil
}
