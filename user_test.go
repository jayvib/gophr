package main

import "testing"

type MockUserStore struct {
	findUser        *User
	findEmailUser   *User
	findUsernameUer *User
	saveUser        *User
}

// Test for the errors
func TestUserNoUsername(t *testing.T) {
	_, err := NewUser("", "user@example.com", "password")
	if err != errNoUsername {
		t.Error("Expected err to be errNoUsername")
	}
}

func TestUserNoPassword(t *testing.T) {
	_, err := NewUser("", "user@example.com", "")
	if err != errNoPassword {
		t.Error("Expected err to be errNoPassword")
	}
}
