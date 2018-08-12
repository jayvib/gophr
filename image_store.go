package main

import (
	"database/sql"
	"errors"
)

var globalImageStore ImageStore

type DBImageStore struct {
	db *sql.DB
}

func NewDBImageStore() ImageStore {
	return &DBImageStore{
		db: globalMySQLDB,
	}
}

func (store *DBImageStore) Save(image *Image) error {
	_, err := store.db.Exec(
		`
		REPLACE INTO images
			(id, user_id, name, location, description, size, created_at)
		VALUES
			(?, ?, ?, ?, ?, ?, ?)		
		`,
		image.ID,
		image.UserID,
		image.Name,
		image.Location,
		image.Description,
		image.Size,
		image.CreatedAt,
	)
	return err
}

func (store *DBImageStore) Find(id string) (*Image, error) {
	row := store.db.QueryRow(
		`
		SELECT id, user_id, name, location, description, size, created_at
		FROM images
		WHERE id = ?`,
		id,
	)
	image := Image{}
	err := row.Scan(
		&image.ID,
		&image.UserID,
		&image.Name,
		&image.Location,
		&image.Description,
		&image.Size,
		&image.CreatedAt,
	)
	return &image, err
}

func (store *DBImageStore) FindAll(offset int) ([]Image, error) {
	return nil, errors.New("not implemented yet")
}

func (store *DBImageStore) FindAllByUser(user *User, offset int) ([]Image, error) {
	return nil, errors.New("not implemented yet")
}
