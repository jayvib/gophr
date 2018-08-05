package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

var (
	errNotImplemented = errors.New("not implemented yet")
)

var globalSessionStore SessionStore

func init() {
	sessionStore, err := NewFileSessionStore("./data/sessions.json")
	if err != nil {
		panic(err)
	}
	globalSessionStore = sessionStore
}

// SessionStore is a generic session storage which can be implemented
// in various ways.
type SessionStore interface {
	Find(string) (*Session, error)
	Save(*Session) error
	Delete(*Session) error
}

func NewFileSessionStore(filename string) (*FileSessionStore, error) {
	store := &FileSessionStore{
		filename: filename,
		Sessions: make(map[string]*Session),
	}
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}
	err = json.Unmarshal(contents, store)
	if err != nil {
		return nil, err
	}
	return store, nil
}

type FileSessionStore struct {
	filename string
	Sessions map[string]*Session
}

func (s *FileSessionStore) Find(sessId string) (*Session, error) {
	if sessId == "" {
		return nil, errors.New("session id can't be empty")
	}
	session, ok := s.Sessions[sessId]
	if !ok {
		return nil, errors.New("session not found")
	}
	return session, nil
}

func (s *FileSessionStore) Save(session *Session) error {
	s.Sessions[session.ID] = session
	contents, err := json.MarshalIndent(s.Sessions, "", "	")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.filename, contents, 0666)
}

func (s *FileSessionStore) Delete(session *Session) error {
	delete(s.Sessions, session.ID)
	contents, err := json.MarshalIndent(s.Sessions, "", "	")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.filename, contents, 0666)
}
