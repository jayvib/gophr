package main

import (
	"fmt"
	"os"
	"flag"
	"net/http"
	"log"
)

var (
	commitID string
)

var (
	PORT string
)

func init() {
	// For printing the version and the commit of the app.
	if len(os.Args) == 2 {
		if os.Args[1] == "version" {
			fmt.Printf("gophr %s\n", commitID)
		}
	}

	flag.StringVar(&PORT, "port", "8080", "The port of gophr app.")
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		RenderTemplate(w, r, "index/home", nil)
	}) // test drive.
	mux.Handle("/assets/",
		http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	log.Printf("Serving gophr at port %s\n", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), mux)
	if err != nil {
		log.Fatal(err.Error())
	}
}