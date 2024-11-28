package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"todo/database"
)

type Task struct {
	Id   int    `json:"Id" db:"id"`
	Desc string `json:"Desc" db:"description"`
}

// db variable to store the database
// var db *sql.DB

// func SetDB(database *sql.DB) {
// 	db = database
// }

func Add(w http.ResponseWriter, r *http.Request) {

	//request
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)

	if err != nil || newTask.Desc == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	//extracting session
	cookie, _ := r.Cookie("session_id")

	// extract user from db using cookie
	username := ""
	err = database.TODO.QueryRow("select username from session where session_id=$1", cookie.Value).Scan(&username)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// id selection
	err = database.TODO.QueryRow(`SELECT
    CASE
        WHEN (SELECT id FROM tasks WHERE id = 1 AND username = $1) IS NULL THEN 1
        ELSE
            (
                SELECT COALESCE(MIN(t1.id + 1), 1)
                FROM tasks t1
                LEFT JOIN tasks t2 ON t1.id + 1 = t2.id AND t1.username = t2.username
                WHERE t2.id IS NULL
                AND t1.username = $1
            )
    END `, username).Scan(&newTask.Id)
	if err != nil {
		http.Error(w, "Error generating ID", http.StatusInternalServerError)
		return
	}

	//insertion
	_, err = database.TODO.Exec("INSERT INTO tasks (id,description,username) VALUES ($1, $2, $3)", newTask.Id, newTask.Desc, username)
	if err != nil {
		http.Error(w, "Error inserting task", http.StatusInternalServerError)
		return
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Task added successfully!",
		"task":    newTask,
	}
	json.NewEncoder(w).Encode(response)

}

func List(w http.ResponseWriter, r *http.Request) {

	// extracting session id
	cookie, _ := r.Cookie("session_id")

	//fetching data
	rows, err := database.TODO.Query("select t.id, t.description FROM tasks t join session s on t.username=s.username where s.session_id = $1", cookie.Value)
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Desc); err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	if len(tasks) == 0 {
		json.NewEncoder(w).Encode(map[string]string{"message": "No Task Found"})
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func Update(w http.ResponseWriter, r *http.Request) {

	//extracting id from body
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	id := newTask.Id
	if err != nil || id <= 0 || newTask.Desc == "" {
		http.Error(w, "Invalid task ID or description", http.StatusBadRequest)
		return
	}

	//extracting session
	cookie, _ := r.Cookie("session_id")

	// extract user from db using cookie
	username := ""
	err = database.TODO.QueryRow("select username from session where session_id=$1", cookie.Value).Scan(&username)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	//updating the task
	var result sql.Result
	result, err = database.TODO.Exec("UPDATE tasks SET description = $2 WHERE id = $1 and username = $3", newTask.Id, newTask.Desc, username)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	//get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error getting rows affected", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Task updated successfully!",
		"task":    newTask,
	})

}

func Delete(w http.ResponseWriter, r *http.Request) {

	//extracting id from body
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	id := newTask.Id

	if err != nil || id <= 0 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	//extracting session
	cookie, _ := r.Cookie("session_id")

	// extract user from db using cookie
	username := ""
	err = database.TODO.QueryRow("select username from session where session_id=$1", cookie.Value).Scan(&username)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	//removing the task

	var result sql.Result
	result, err = database.TODO.Exec("DELETE FROM tasks WHERE id = $1 and username = $2", id, username)
	if err != nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}

	//get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error getting rows affected", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})

}
