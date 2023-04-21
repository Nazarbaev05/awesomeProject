package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func main() {
	// Set up database connection
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Define routes and their handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			name := r.FormValue("name")
			email := r.FormValue("email")
			password := r.FormValue("password")

			// Insert user data into database
			result, err := db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", name, email, password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			id, _ := result.LastInsertId()

			data := map[string]interface{}{"id": id, "name": name, "email": email, "password": password}

			// Redirect to profile page with user ID
			t, _ := template.ParseFiles("profile.html")
			t.Execute(w, data)
			return
		}
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
	})

	// Start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
