package main

import (
	"database/sql"
	"fmt"
	//_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func main() {

	http.Handle("/", http.FileServer(http.Dir("./static")))
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			name := r.FormValue("name")
			email := r.FormValue("email")
			password := r.FormValue("password")

			result, err := db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", name, email, password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			id, _ := result.LastInsertId()

			data := map[string]interface{}{"id": id, "name": name, "email": email, "password": password}

			t, _ := template.ParseFiles("profile.html")
			t.Execute(w, data)
			return
		}
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
	})
	fmt.Println("Server is running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
