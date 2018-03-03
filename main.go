package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alecthomas/template"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
)

// StaticInfo ...
type StaticInfo struct {
	NewsletterName string
	PublisherURL   string
	PublisherName  string
}

var (
	siteName            = "Dart Weekly Newsletter"
	siteBlurb           = "All the latest news on the Dart programming language"
	siteAddress         = "https://testdomain.com"
	localAddressAndPort = "127.0.0.1:8000"

	// Templates
	landingTemplate  = template.Must(template.ParseFiles("templates/landing.tmpl"))
	confirmTemplate  = template.Must(template.ParseFiles("templates/confirm.tmpl"))
	thankyouTemplate = template.Must(template.ParseFiles("templates/thankyou.tmpl"))
	privacyTemplate  = template.Must(template.ParseFiles("templates/privacy.tmpl"))
	issuesTemplate   = template.Must(template.ParseFiles("templates/issues.tmpl"))
	archivesTemplate = template.Must(template.ParseFiles("templates/archives.tmpl"))
	latestTemplate   = template.Must(template.ParseFiles("templates/latest.tmpl"))

	newsletterInfo = StaticInfo{
		NewsletterName: "Dart Weekly",
		PublisherURL:   "https://saltypress.com",
		PublisherName:  "SaltyPress",
	}
)

func landingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)
	w.WriteHeader(http.StatusOK)
	data := struct {
		StaticInfo      StaticInfo
		NewsletterBlurb string
	}{
		StaticInfo:      newsletterInfo,
		NewsletterBlurb: siteBlurb,
	}
	fmt.Printf("%#+v", data)
	landingTemplate.Execute(w, data)
}

// ConfirmationHandler is <siteAddress>/confirm/TOKEN
func confirmationHandler(w http.ResponseWriter, r *http.Request) {
	// Confirm in DB and give them a confirmation.html
	vars := mux.Vars(r)

	log.Println(vars)
	log.Println(vars["confirmationID"])
	w.WriteHeader(http.StatusOK)
}

func thanksHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)
	w.WriteHeader(http.StatusOK)
	data := struct {
		StaticInfo StaticInfo
	}{
		StaticInfo: newsletterInfo,
	}
	thankyouTemplate.Execute(w, data)
}

func privacyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)
	w.WriteHeader(http.StatusOK)
	data := struct {
		StaticInfo StaticInfo
	}{
		StaticInfo: newsletterInfo,
	}
	privacyTemplate.Execute(w, data)
}

func issueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var issue string
	var ok bool
	if issue, ok = vars["issue"]; !ok {
		log.Println("Not ok, no val")
		// TODO: Get latest
		return
	}
	w.WriteHeader(http.StatusOK)
	issuesTemplate.Execute(w, issue)
}

func latestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	data := struct {
		StaticInfo StaticInfo
	}{
		StaticInfo: newsletterInfo,
	}
	latestTemplate.Execute(w, data)
}

func archivesHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	archivesTemplate.Execute(w, struct{}{})
}

func main() {

	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	database, _ := sql.Open("sqlite3", "./content.db")
	statement, _ := database.Prepare(
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY,
			issue_number INTEGER,
			post_title TEXT, 
			post_link TEXT,
			description TEXT,
			post_date TEXT
		)`,
	)
	res, err := statement.Exec()
	log.Println("Created Table result:", res, err)

	// Add demo post
	_, err = database.Exec(
		`INSERT INTO posts
		(
			issue_number,
			post_title,
			post_link,
			description,
			post_date
		)
		VALUES
		(?,?,?,?,?)
		`,
		1, // First issue
		"Hairy Balls",
		"http://www.spaceship.com.au",
		"This is the description",
		SqlLiteDate(time.Now()),
	)
	if err != nil {
		log.Println("Insert failed", err)
	}
	rows, err := database.Query(`
		SELECT 
			issue_number,
			post_title,
			post_link,
			description,
			post_date
		FROM posts`,
	)
	if err != nil {
		log.Println("uh oh:", err)
	}
	var (
		issueNumber int
		postTitle   string
		postLink    string
		description string
		postDate    SqlLiteDate
	)
	for rows.Next() {
		err = rows.Scan(
			&issueNumber,
			&postTitle,
			&postLink,
			&description,
			&postDate,
		)
		if err != nil {
			log.Println("Uh oh", err)
			continue
		}
		fmt.Println(
			issueNumber,
			postTitle,
			postLink,
			description,
			postDate,
		)
	}

	redirectHome := http.RedirectHandler("/", 307)

	r := mux.NewRouter()
	// Landing Page
	r.HandleFunc("/", landingHandler)
	// Privacy Page
	r.HandleFunc("/privacy", privacyHandler)
	r.HandleFunc("/confirm/{confirmationID}", confirmationHandler) // takes a confirmationID

	r.HandleFunc("/subscribe", thanksHandler).Methods("POST")
	r.Handle("/subscribe", redirectHome).Methods("GET")

	r.HandleFunc("/latest", latestHandler)
	r.HandleFunc("/issues/{issue:[0-9]+}", issueHandler)
	r.HandleFunc("/archives", archivesHandler)

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
