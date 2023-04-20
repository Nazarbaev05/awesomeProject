package awesomeProject

import (
	"database/sql"
	"fmt"
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
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/database")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Define routes and their handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Serve registration form
			tmpl, err := template.ParseFiles("index.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
		} else if r.Method == "POST" {
			// Parse form data
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

			// Redirect to profile page with user ID
			http.Redirect(w, r, fmt.Sprintf("/profile/%d", id), http.StatusSeeOther)
		}
	})

	http.HandleFunc("/profile/", func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from URL parameter
		id := r.URL.Path[len("/profile/"):]

		// Query database for user data
		var user User
		err := db.QueryRow("SELECT id, name, email, password FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Render profile template with user data
		tmpl, err := template.ParseFiles("profile.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, user)
	})

	// Start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
