package clients

import (
	"testing"
)

//The purpose of this test is to ensure you canreply on the mocks

const TOKEN_MOCK = "this is a token"

func TestSeagullMock_GetCollection(t *testing.T) {

	sc := NewSeagullMock()
	var col struct{ Something string }
	sc.SetMockNextCollectionCall("123.456stuff", `{"Something":"anit no thing"}`, nil)

	err := sc.GetCollection("123.456", "stuff", TOKEN_MOCK, &col)
	if err != nil {
		t.Error("Should not return an error")
	}
	if col.Something != "anit no thing" {
		t.Error("Should have given mocked collection")
	}

	errMsg := "Testing my error"
	sc.SetMockNextCollectionCall("123.456stuff", "", errors.New(errMsg))

	err = sc.GetCollection("123.456", "stuff", TOKEN_MOCK, &col)
	if err == nil {
		t.Error("Should return an error")
	}
	if err.Error() != errMsg {
		t.Errorf("Unexepected error message expected:%v/receieved:%v", errMsg, err.Error())
	}
}

func TestSeagullMock_GetPrivatePair(t *testing.T) {
	sc := NewSeagullMock()

	if pp := sc.GetPrivatePair("123.456", "Stuff", TOKEN_MOCK); pp == nil {
		t.Error("Should give us mocked private pair")
	}

}
