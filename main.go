package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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

	router := NewRouter()
	router.Handle("GET", "/", HandleHome)
	router.Handle("GET", "/register", HandleNewUserPage) // user to display the registration page
	router.Handle("POST", "/register", HandleCreateUser) // use to register the user to the system
	router.Handle("GET", "/login", HandleNewSessionPage)
	router.Handle("POST", "/login", HandleSessionCreate)
	router.Handle("GET", "/sign-out", HandleSessionDestroy)
	router.Handle("GET", "/account", HandleUserEdit)
	router.Handle("POST", "/account", HandleUserUpdate)
	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	// secureRouter := NewRouter()
	// secureRouter.Handle("GET", "/sign-out", HandleSessionDestroy)
	// securedMiddleware := NewMiddleware(
	// 	http.HandleFunc(RequireLogin),
	// )

	log.Printf("Serving gophr at port %s\n", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), router)
	if err != nil {
		log.Fatal(err.Error())
	}
}
