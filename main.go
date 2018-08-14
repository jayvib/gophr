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

	// User Store
	store, err := NewFileUserStore("./data/users.json")
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating user store: %s", err))
	}
	globalUserStore = store

	// Session Store
	sessionStore, err := NewFileSessionStore("./data/sessions.json")
	if err != nil {
		panic(err)
	}
	globalSessionStore = sessionStore

	db, err := NewMySQLDB("root:mysql123@tcp(127.0.0.1:3306)/gophr")
	if err != nil {
		panic(err)
	}
	globalMySQLDB = db

	globalImageStore = NewDBImageStore()
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
	router.Handle("GET", "/images/new", HandleImageNew)
	router.Handle("POST", "/images/new", HandleImageCreate)
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
