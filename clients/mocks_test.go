package clients

import (
	"errors"
	"reflect"
	"testing"
)

//The purpose of this test is to ensure you canreply on the mocks

const USERID, GROUPID, TOKEN_MOCK = "123user", "456group", "this is a token"

func makeExpectedPermissons() Permissions {
	return Permissions{"root": Allowed}
}

func makeExpectedUsersPermissions() UsersPermissions {
	return UsersPermissions{GROUPID: Permissions{GROUPID: Allowed}}
}

func TestGatekeeperMock_UserInGroup(t *testing.T) {

	expected := makeExpectedPermissons()

	gkc := NewGatekeeperMock(nil, nil)

	if perms, err := gkc.UserInGroup(USERID, GROUPID); err != nil {
		t.Fatal("No error should be returned")
	} else if !reflect.DeepEqual(perms, expected) {
		t.Fatalf("Perms where [%v] but expected [%v]", perms, expected)
	}
}

func TestGatekeeperMock_UsersInGroup(t *testing.T) {

	expected := makeExpectedUsersPermissions()

	gkc := NewGatekeeperMock(nil, nil)

	if perms, err := gkc.UsersInGroup(GROUPID); err != nil {
		t.Fatal("No error should be returned")
	} else if !reflect.DeepEqual(perms, expected) {
		t.Fatalf("Perms were [%v] but expected [%v]", perms, expected)
	}
}

func TestGatekeeperMock_GroupsForUser(t *testing.T) {
	expected := makeExpectedUsersPermissions()

	gkc := NewGatekeeperMock(nil, nil)
	if perms, err := gkc.GroupsForUser(GROUPID); err != nil {
		t.Fatal("No error should be returned")
	} else if !reflect.DeepEqual(perms, expected) {
		t.Fatalf("Perms were [%v] but expected [%v]", perms, expected)
	}

	// testing with error
	gkc.SetExpected(nil, errors.New("gk error"))
	if _, err := gkc.GroupsForUser(GROUPID); err == nil {
		t.Fatal("An error should be returned")
	}

	// testing with permisson set
	mockedPerms := Permissions{"view": Allowed, "root": Allowed}
	gkc.SetExpected(mockedPerms, nil)
	expectedPerms := UsersPermissions{}
	expectedPerms[GROUPID] = mockedPerms
	if perms, err := gkc.GroupsForUser(GROUPID); err != nil {
		t.Fatal("No error should be returned")
	} else if !reflect.DeepEqual(perms, expectedPerms) {
		t.Fatalf("Perms were [%v] but expected [%v]", perms, expectedPerms)
	}

	//testing with UserIds
	gkc.SetExpected(mockedPerms, nil)
	gkc.UserIDs = []string{"user1", "user2"}
	expectedPerms = UsersPermissions{}
	for _, user := range gkc.UserIDs {
		expectedPerms[user] = mockedPerms
	}
	expectedPerms[GROUPID] = Permissions{"root": Allowed}
	if userPerms, err := gkc.GroupsForUser(GROUPID); err != nil {
		t.Fatal("No error should be returned")
	} else if !reflect.DeepEqual(userPerms, expectedPerms) {
		t.Fatalf("Perms were [%v] but expected [%v]", userPerms, expectedPerms)
	}

}
func TestGatekeeperMock_SetPermissions(t *testing.T) {

	gkc := NewGatekeeperMock(nil, nil)

	expected := makeExpectedPermissons()

	if perms, err := gkc.SetPermissions(USERID, GROUPID, expected); err != nil {
		t.Fatal("No error should be returned")
	} else if !reflect.DeepEqual(perms, expected) {
		t.Fatalf("Perms where [%v] but expected [%v]", perms, expected)

	}
}

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
