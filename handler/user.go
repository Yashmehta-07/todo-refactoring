package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "fmt"
	"net/http"
	"time"
	"todo/database"
	dbhelper "todo/database/dbHelper"

	_ "github.com/lib/pq" // Import pq driver
	// "fmt"
)

// User
type User struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

// Register
func Register(w http.ResponseWriter, r *http.Request) {

	//request
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	//insertion
	_, err = database.TODO.Exec("INSERT INTO auth (username,password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		http.Error(w, "Error inserting task or user already exists", http.StatusInternalServerError)
		return
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "success",
	})

}

// Login
func Login(w http.ResponseWriter, r *http.Request) {

	//request
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	//fetching data
	var username string
	err = database.TODO.Get(&username, "SELECT username FROM auth WHERE username = $1 AND password = $2", user.Username, user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}

	//generating session
	session_id, err := dbhelper.GenerateSessionID()
	if err != nil {
		http.Error(w, "Error generating session", http.StatusInternalServerError)
		return
	}

	fmt.Print(session_id)

	//insertion
	_, err = database.TODO.Exec("INSERT INTO session (session_id,username,created_at) VALUES ($1, $2, $3)", session_id, user.Username, time.Now().UTC())
	if err != nil {
		http.Error(w, "Error inserting task", http.StatusInternalServerError)
		return
	}

	//set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session_id,
		HttpOnly: true,
		Secure:   true,
		// Expires:  time.Now().Add(2 * time.Minute)
		SameSite: http.SameSiteLaxMode,
	})

	//response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "login successfull",
	})

}

// Logout

func Logout(w http.ResponseWriter, r *http.Request) {

	// session check
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		http.Error(w, "already logout ", http.StatusUnauthorized)
		return
	}

	// deleting cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",    // Name of the cookie to delete
		Value:    "",              // Empty the value
		Expires:  time.Unix(0, 0), // Expire in the past
		MaxAge:   -1,              // Invalidate immediately
		HttpOnly: true,            // Keep HttpOnly for security
		Secure:   true,
	})

	//deleting session
	_, err = database.TODO.Exec("DELETE FROM session WHERE session_id = $1", cookie.Value)
	if err != nil {
		http.Error(w, "Error deleting session", http.StatusInternalServerError)
		return
	}

	if r.URL.Path == "/logout" {
		//response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "logout successfull",
		})
	}
}
