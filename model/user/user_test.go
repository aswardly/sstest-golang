//user_test provides unit tests for business domain model of user
package user_test

import (
	"fmt"
	"os"
	"sstest/model/user"
	"testing"
)

var newUser, activeUser, inactiveUser, suspendedUser *user.User

func TestMain(m *testing.M) {
	//test setup
	newUser, _ = user.New("newUser", "New User", "New User Address")

	activeUser, _ = user.New("activeUser", "Active User", "Active User Address")
	activeUser.SetStatus(user.StatusActive)

	inactiveUser, _ = user.New("inactiveUser", "Inactive User", "Inactive User Address")
	inactiveUser.SetStatus(user.StatusInactive)

	suspendedUser, _ = user.New("suspendedUser", "Suspended User", "Suspended User Address")
	suspendedUser.SetStatus(user.StatusSuspended)

	//run tests
	exitCode := m.Run()
	//test teardown
	os.Exit(exitCode)
}

func TestNewlyCreatedUser(t *testing.T) {
	t.Run("Status Must Be Inactive", func(t *testing.T) {
		if user.StatusInactive != newUser.Status() {
			t.Errorf("expected %v but got %v", user.StatusInactive, newUser.Status())
		}
	})
	t.Run("Default Password Must Be changeme", func(t *testing.T) {
		match, err := newUser.ValidatePassword("changeme")
		if err != nil {
			t.Errorf("validating password got error %v", err)
		}
		if false == match {
			t.Errorf("validating password got response %v", match)
		}
	})
}

func TestSettingPassword(t *testing.T) {
	var match bool
	if _, err := newUser.SetPassword("testPassword"); err != nil {
		t.Errorf("setting new password got error %v", err)
	}
	match, err := newUser.ValidatePassword("testPassword")
	if err != nil {
		t.Errorf("validating password got error %v", err)
	}
	if false == match {
		t.Error("password does not match\n")
	}
}

func TestSettingInvalidPasswordHash(t *testing.T) {
	if _, err := newUser.SetPasswordHash([]byte("anInvalidPasswordHash")); err == nil {
		t.Error("expected error but got none\n")
	}
}

func TestSettingUnknownStatusValue(t *testing.T) {
	_, err := newUser.SetStatus("UnknownStatus")
	if err == nil {
		t.Error("expected error but got none\n")
	}
}

func TestCanOrder(t *testing.T) {
	activeUserOrder, _ := activeUser.CanOrder()
	inactiveUserOrder, _ := inactiveUser.CanOrder()
	suspendedUserOrder, _ := suspendedUser.CanOrder()

	var canOrderTests = []struct {
		testCase string
		expected bool
		actual   bool
	}{
		{"Active User Order", true, activeUserOrder},
		{"Inactive User Order", false, inactiveUserOrder},
		{"Suspended User Order", false, suspendedUserOrder},
	}

	for _, test := range canOrderTests {
		t.Run(fmt.Sprintf("%s", test.testCase), func(t *testing.T) {
			if test.expected != test.actual {
				t.Errorf("want %v got %v", test.expected, test.actual)
			}
		})
	}
}
