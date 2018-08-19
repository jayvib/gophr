package main

import (
	"fmt"
	"net/http"
)

// HandleHome is the callback to display the home page of gophr
func HandleHome(w http.ResponseWriter, r *http.Request) {
	images, err := globalImageStore.FindAll(0)
	if err != nil {
		panic(err)
	}
	fmt.Println("Number of Images found: ", len(images))
	for _, image := range images {
		fmt.Println(image.ShowRoute(), image.StaticRoute())
	}
	RenderTemplate(w, r, "index/home", map[string]interface{}{
		"Images": images,
	})
}
