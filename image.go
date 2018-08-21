package main

import (
	"fmt"
	"image"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/disintegration/imaging"
)

const (
	imageIDLength  = 10
	widthThumbnail = 400
	widthPreview   = 800
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func NewImage(user *User) *Image {
	return &Image{
		ID:        GenerateID("img", imageIDLength),
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}
}

var mimeExtensions = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
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

func (i *Image) StaticRoute() string {
	return "/im/" + i.Location
}

func (i *Image) ShowRoute() string {
	return "/image/" + i.ID
}

// CreateResizedImages resizes the image to preview and thumbnail size.
func (i *Image) CreateResizedImages() error {
	srcImage, err := imaging.Open("./data/images/" + i.Location)
	if err != nil {
		return err
	}

	errChan := make(chan error)
	go i.resizePreview(errChan, srcImage)
	go i.resizeThumbnail(errChan, srcImage)

	for i := 0; i < 2; i++ {
		e := <-errChan
		if e == nil {
			err = e
		}
	}

	return err
}

// resizePreview resizes the image to be use for preview
func (i *Image) resizePreview(errorChan chan error, srcImage image.Image) {
	fmt.Println("Resizing for preview image")
	size := srcImage.Bounds().Size()
	ratio := float64(size.Y) / float64(size.X)
	targetHeight := int(float64(widthPreview) * ratio)
	dstImage := imaging.Resize(srcImage, widthPreview, targetHeight, imaging.Lanczos)
	destination := "./data/images/preview/" + i.Location
	errorChan <- imaging.Save(dstImage, destination)

}

// resizeThumbnail resize the image to be use for thumbnails
func (i *Image) resizeThumbnail(errorChan chan error, srcImage image.Image) {
	fmt.Println("Resizing fore thumbnail image")
	dstImage := imaging.Thumbnail(srcImage, widthThumbnail, widthThumbnail, imaging.Lanczos)
	destination := "./data/images/thumbnail/" + i.Location
	errorChan <- imaging.Save(dstImage, destination)
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
	err = i.CreateResizedImages()
	if err != nil {
		return err
	}
	return globalImageStore.Save(i)
}

// CreateFromFile get the file uploaded by the user in a multipart form and save the file info
// to the database and file to assets/data/images.
func (i *Image) CreateFromFile(file multipart.File, headers *multipart.FileHeader) error {
	i.Name = headers.Filename
	i.Location = i.ID + filepath.Ext(i.Name)
	savedFile, err := os.Create("./data/images/" + i.Location)
	if err != nil {
		return err
	}
	defer savedFile.Close()
	size, err := io.Copy(savedFile, file)
	if err != nil {
		return err
	}
	i.Size = size
	err = i.CreateResizedImages()
	if err != nil {
		return err
	}
	return globalImageStore.Save(i)
}
