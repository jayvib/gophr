package main

import "testing"

func TestSessionStore(t *testing.T) {
	var sessionStore SessionStore

	t.Run("File Session Store", func(t *testing.T) {
		fileSessionStore := &FileSessionStore{
			filename: "session.json",
			Sessions: make(map[string]*Session),
		}
		sessionStore = fileSessionStore

		t.Run("File Session Store - Find", func(t *testing.T) {
			// Find session
			sess, err := sessionStore.Find("")
			if err == nil {
				t.Error("expecting an error but haven't got one")
			}
			if sess != nil {
				t.Error("empty session id must be returned nil")
			}

			sess, err = sessionStore.Find("sess_1234567890")
			if err != nil {
				t.Error(err.Error())
			}

			if sess == nil {
				t.Error("expecting not empty session")
			}
		})

		t.Run("File Session Store - Save", func(t *testing.T) {
			t.Skip("can't be test yet unless the find method is already passed")
			// Test if the the session is nil it will return an error
			session := &Session{
				ID:     "sess_test123",
				UserID: "test123",
			}
			err := sessionStore.Save(nil)
			if err == nil {
				t.Error("Expecting an error to be receive on saving nil session")
			}
			// Test if the user already exist it will return an error
			err = sessionStore.Save(session)
			if err != nil {
				t.Error(err)
			}

			storedSess, err := sessionStore.Find(session.ID)
			if err != nil {
				t.Error(err)
			}
			if storedSess.UserID != "test123" {
				t.Errorf("expecting user ID %s but got %s\n", "test123", storedSess.UserID)
			}
		})
	})
}
