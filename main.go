package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("take-me-out-of-code"))
var templates = template.Must(template.ParseFiles("views/user/login.html", "views/user/create.html"))

func main() {

	var port = flag.String("port", ":8080", "network port to receive http requests over")

	flag.Parse()

	router := mux.NewRouter()

	router.HandleFunc("/", index)
	router.HandleFunc("/user/login", loginView)
	router.HandleFunc("/user/create", createView).Methods("GET")
	router.HandleFunc("/user/create", createAction).Methods("POST")

	http.Handle("/", router)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {

	session, _ := store.Get(req, "user")

	if _, exists := session.Values["user_id"]; !exists {
		http.Redirect(w, req, "/user/login", 302)
		return
	}

	fmt.Fprintln(w, "This does nothing :(")
}

func loginView(w http.ResponseWriter, req *http.Request) {

	session, _ := store.Get(req, "user")

	data := ViewData{"", nil}
	if flashes := session.Flashes("feedback"); len(flashes) > 0 {
		data.Feedback = flashes[0].(string)
		session.Save(req, w)
	}

	err := templates.ExecuteTemplate(w, "login.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createView(w http.ResponseWriter, req *http.Request) {

	err := templates.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createAction(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "user")

	//todo: validate
	//todo: create account

	session.AddFlash("Successfully Created Account", "feedback")
	session.Save(req, w)
	http.Redirect(w, req, "/user/login", 302)
}
