package clients

import (
	"encoding/json"
)

type (
	GatekeeperMock struct {
		expectedPermissions Permissions
		expectedError       error
		UserIDs             []string
	}
	SeagullMock struct{}
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
	return &SeagullMock{}
}

func (mock *SeagullMock) GetPrivatePair(userID, hashName, token string) *PrivatePair {
	return &PrivatePair{ID: "mock.id.123", Value: "mock value"}
}

func (mock *SeagullMock) GetCollection(userID, collectionName, token string, v interface{}) error {
	json.Unmarshal([]byte(`{"Something":"anit no thing", "patient": {"birthday": "2016-01-01"}}`), &v)
	return nil
}
