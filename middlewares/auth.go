package middlewares

import (
	"database/sql"
	"net/http"
	"time"
	"todo/database"
	"todo/handler"
	"todo/logging"
)

// middlewares
func Caller(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//session check
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			logging.Log(err, "Unauthorized", "warning", 401, r)
			return
		}
		sessionID := cookie.Value

		//fetching data
		data := struct {
			Username   string    `db:"username"`
			Created_at time.Time `db:"created_at"`
		}{}

		err = database.TODO.Get(&data, "SELECT username, created_at FROM session WHERE session_id = $1", sessionID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				logging.Log(err, "User not found", "warning", 404, r)
				return
			}
			http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
			logging.Log(err, "Error fetching tasks", "warning", 500, r)
			return
		}

		duration := time.Now().UTC().Sub(data.Created_at) //time.Since(created_at)
		if duration >= 1*time.Hour {
			handler.Logout(w, r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			logging.Log(err, "Unauthorized", "warning", 401, r)
			return

		}

		next.ServeHTTP(w, r)
	})
}
