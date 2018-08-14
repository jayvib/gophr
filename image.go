package main

import (
	"time"
	"net/http"
	"mime"
	"path/filepath"
	"os"
	"io"
	"mime/multipart"
)

const imageIDLength = 10

func NewImage(user *User) *Image {
	return &Image{
		ID: GenerateID("img", imageIDLength),
		UserID: user.ID,
		CreatedAt: time.Now(),
	}
}

var mimeExtensions = map[string]string{
	"image/png": ".png",
	"image/jpeg": ".jpg",
	"image/gif": ".gif",
}

type Image struct {
	ID          string
	UserID      string
	Name        string
	Location    string
	Size        int64
	CreatedAt   time.Time
	Description string
}

type ImageStore interface {
	Save(image *Image) error
	Find(id string) (*Image, error)
	FindAll(offset int) ([]Image, error)
	FindAllByUser(user *User, offset int) ([]Image, error)
}

func (i *Image) CreateFromURL(imageURL string) error {
	response, err := http.Get(imageURL)
	if err != nil {
		return errImageURLInvalid
	}
	if response.StatusCode != http.StatusOK {
		return errImageURLInvalid
	}
	defer response.Body.Close()
	mimeType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return errInvalidImageType
	}
	ext, ok := mimeExtensions[mimeType]
	if !ok {
		return errInvalidImageType
	}
	i.Name = filepath.Base(imageURL)
	i.Location = i.ID + ext
	savedFile, err := os.Create("./data/images/" + i.Location)
	if err != nil {
		return err
	}
	defer savedFile.Close()
	size, err := io.Copy(savedFile, response.Body)
	if err != nil {
		return err
	}
	i.Size = size
	return globalImageStore.Save(i)
}

func (i *Image) CreateFromFile(file multipart.File, headers *multipart.FileHeader) error {
	i.Name = headers.Filename
	i.Location = i.ID + filepath.Ext(i.Name)
	savedFile, err := os.Create("./data/images" + i.Location)
	if err != nil {
		return err
	}
	defer savedFile.Close()
	size, err := io.Copy(savedFile, file)
	if err != nil {
		return err
	}
	i.Size = size
	return globalImageStore.Save(i)
}