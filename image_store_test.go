package main

import "testing"

var iStore ImageStore

func setup() {
	db, err := NewMySQLDB("root:mysql123@tcp(127.0.0.1:3306)/gophr")
	if err != nil {
		panic(err)
	}
	globalMySQLDB = db
	iStore = NewDBImageStore()
}

func TestFindAllByUser(t *testing.T) {
	setup()
	userId := "usr_xQ3luRi0X7"
	user := &User{
		ID: userId,
	}
	images, err := iStore.FindAllByUser(user, 0)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(images) != 2 {
		t.Errorf("expecting 2 images but got %d", len(images))
	}
}
