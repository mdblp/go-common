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
