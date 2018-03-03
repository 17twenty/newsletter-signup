package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alecthomas/template"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
)

func LandingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "LandingHandler", vars)
}

// ConfirmationHandler is <siteAddress>/confirm/TOKEN
func ConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	// Confirm in DB and give them a confirmation.html
	vars := mux.Vars(r)

	log.Println(vars)
	log.Println(vars["confirmationID"])
	w.WriteHeader(http.StatusOK)
}

func ThanksHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "thanks: %v\n", vars["category"])
}

func LatestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "thanks: %v\n", vars["category"])
}

func ArchivesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "thanks: %v\n", vars["category"])
}

func TestTemplateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)
	w.WriteHeader(http.StatusOK)

	data := struct {
		PageTitle string
		Todos     []struct {
			Title string
			Done  bool
		}
	}{
		PageTitle: "My TODO list",
		Todos: []struct {
			Title string
			Done  bool
		}{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}

	testTemplate.Execute(w, data)
}

var (
	siteAddress         = "https://testdomain.com"
	localAddressAndPort = "127.0.0.1:8000"

	// Templates

	testTemplate = template.Must(template.ParseFiles("templates/testing.html"))
)

func main() {

	database, _ := sql.Open("sqlite3", "./content.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	statement.Exec("Demo", "Person")
	rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
	var id int
	var firstname string
	var lastname string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", LandingHandler)
	r.HandleFunc("/confirm/{confirmationID}", ConfirmationHandler) // takes a confirmationID
	r.HandleFunc("/thankyou", ThanksHandler).Methods("POST")
	r.HandleFunc("/latest", LatestHandler)
	r.HandleFunc("/archives", ArchivesHandler)
	r.HandleFunc("/testing", TestTemplateHandler)

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	srv := &http.Server{
		Handler: r,
		Addr:    localAddressAndPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
