package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	publicRouter := router.PathPrefix("/").Subrouter()

	// TODO: Find a way to separate the routes that need authentication
	publicRouter.HandleFunc("/sign-out", AuthMiddleware(HandleSessionDestroy)).Methods("GET")
	publicRouter.HandleFunc("/account", AuthMiddleware(HandleUserEdit)).Methods("GET")
	publicRouter.HandleFunc("/account", AuthMiddleware(HandleUserUpdate)).Methods("POST")
	publicRouter.HandleFunc("/images/new", AuthMiddleware(HandleImageNew)).Methods("GET")
	publicRouter.HandleFunc("/images/new", AuthMiddleware(HandleImageCreate)).Methods("POST")

	publicRouter.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	publicRouter.HandleFunc("/register", HandleNewUserPage).Methods("GET") // user to display the registration page
	publicRouter.HandleFunc("/register", HandleCreateUser).Methods("POST") // use to register the user to the system
	publicRouter.HandleFunc("/", HandleHome).Methods("GET")
	publicRouter.HandleFunc("/login", HandleNewSessionPage).Methods("GET")
	publicRouter.HandleFunc("/login", HandleSessionCreate).Methods("POST")

	router.Use(loggingMiddleware)
	log.Printf("Serving gophr at port %s\n", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), router)
	if err != nil {
		log.Fatal(err.Error())
	}
}
