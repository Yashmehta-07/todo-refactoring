package handler

import (
	"database/sql"
	"encoding/json"
	_ "fmt"
	"net/http"
	"time"
	"todo/database"
	dbhelper "todo/database/dbHelper"
	"todo/logging"

	_ "github.com/lib/pq" // Import pq driver
	// "fmt"
)

// User
type User struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user with a username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "User registration data"
// @Success 200 {object} map[string]interface{} "Registration successful"
// @Failure 400 {object} map[string]string "Invalid username or password"
// @Failure 500 {object} map[string]string "Error inserting user or user already exists"
// @Router /register [post]
func Register(w http.ResponseWriter, r *http.Request) {

	//request
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		logging.Log(err, "Invalid username or password", "warning", 400, r)
		return
	}

	//insertion
	_, err = database.TODO.Exec("INSERT INTO auth (username,password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		http.Error(w, "Error inserting task or user already exists", http.StatusInternalServerError)
		logging.Log(err, "Error inserting task or user already exists", "error", 500, r)
		return
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "success",
	})

	logging.Log(err, "success", "info", 200, r)

}

// Login godoc
// @Summary Login a user
// @Description Authenticate a user and create a session
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "User login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]string "Invalid username or password"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Error logging in"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {

	//request
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		logging.Log(err, "Invalid username or password", "warning", 400, r)
		return
	}

	//fetching data
	var username string
	err = database.TODO.Get(&username, "SELECT username FROM auth WHERE username = $1 AND password = $2", user.Username, user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			logging.Log(err, "User not found", "error", 404, r)
			return
		}
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		logging.Log(err, "Error fetching tasks", "error", 500, r)
		return
	}

	//generating session
	session_id, err := dbhelper.GenerateSessionID()
	if err != nil {
		http.Error(w, "Error generating session", http.StatusInternalServerError)
		logging.Log(err, "Error generating session", "error", 500, r)
		return
	}

	//insertion
	_, err = database.TODO.Exec("INSERT INTO session (session_id,username,created_at) VALUES ($1, $2, $3)", session_id, user.Username, time.Now().UTC())
	if err != nil {
		http.Error(w, "Error inserting task", http.StatusInternalServerError)
		logging.Log(err, "Error inserting task", "error", 500, r)
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

	logging.Log(err, "login successfull", "info", 200, r)

}

// Logout godoc
// @Summary Logout a user
// @Description Logout a user by invalidating their session
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Logout successful"
// @Failure 401 {object} map[string]string "Already logged out or invalid session"
// @Failure 500 {object} map[string]string "Error deleting session"
// @Router /logout [post]
func Logout(w http.ResponseWriter, r *http.Request) {

	// session check
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		http.Error(w, "already logout ", http.StatusUnauthorized)
		logging.Log(err, "already logout ", "warning", 401, r)
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
		logging.Log(err, "Error deleting session", "error", 500, r)
		return
	}

	if r.URL.Path == "/logout" {
		//response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "logout successfull",
		})

		logging.Log(err, "logout successfull", "info", 200, r)
	}

}
