package middlewares

import (
	"database/sql"
	"net/http"
	"time"
	"todo/database"
	"todo/handler"
)

// db variable to store the database
// var db *sql.DB

// func MiddlewareSetDB(database *sql.DB) {
// 	db = database
// }

// middlewares
func Caller(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// db := utils.GetDB()
		// fmt.Print(db)
		//session check
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		sessionID := cookie.Value

		//fetching data
		var (
			username   string
			created_at time.Time
		)
		err = database.TODO.QueryRow("SELECT username, created_at FROM session WHERE session_id = $1", sessionID).Scan(&username, &created_at)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
			return
		}

		duration := time.Now().UTC().Sub(created_at) //time.Since(created_at)
		if duration >= 1*time.Hour {
			handler.Logout(w, r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return

		}

		next.ServeHTTP(w, r)
	})
}