package main

import "errors"

var (
	errNotImplemented = errors.New("not implemented yet")
)

// SessionStore is a generic session storage which can be implemented
// in various ways
type SessionStore interface {
	Find(string) (*Session, error)
	Save(*Session) error
	Delete(*Session) error
}

type FileSessionStore struct {
	filename string
	Sessions map[string]*Session
}

func (s *FileSessionStore) Find(sessId string) (*Session, error) {
	return nil, errNotImplemented
}

func (s *FileSessionStore) Save(session *Session) error {
	return errNotImplemented
}

func (s *FileSessionStore) Delete(session *Session) error {
	return errNotImplemented
}
