package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var globalUserStore UserStore

type UserStore interface {
	Find(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindByUsername(string) (*User, error)
	Save(User) error
}

func init() {
	store, err := NewFileUserStore("./data/users.json")
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating user store: %s", err))
	}
	globalUserStore = store
}

func NewFileUserStore(filename string) (*FileUserStore, error) {
	fileUserStore := &FileUserStore{
		filename: filename,
		Users:    make(map[string]User),
	}

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fileUserStore, nil
		}
		return nil, err
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &fileUserStore.Users)
	if err != nil {
		return nil, err
	}
	return fileUserStore, nil
}

type FileUserStore struct {
	filename string
	Users    map[string]User
}

func (store *FileUserStore) Save(user User) error {
	store.Users[user.ID] = user
	b, err := json.MarshalIndent(store.Users, "", "	")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(store.filename, b, 0660)
	if err != nil {
		return err
	}
	return nil
}

func (store *FileUserStore) FindByEmail(email string) (*User, error) {
	if email == "" {
		return nil, nil
	}
	for _, user := range store.Users {
		if strings.ToLower(user.Email) == strings.ToLower(email) {
			return &user, nil
		}
	}
	return nil, nil
}

func (store *FileUserStore) FindByUsername(username string) (*User, error) {
	if username == "" {
		return nil, nil
	}
	for _, user := range store.Users {
		if strings.ToLower(user.Username) == strings.ToLower(username) {
			return &user, nil
		}
	}
	return nil, nil
}

func (store *FileUserStore) Find(id string) (*User, error) {
	if id == "" {
		return nil, nil
	}
	if user, ok := store.Users[id]; ok {
		return &user, nil
	}
	return nil, nil
}
