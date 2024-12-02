package handler

import (
	"database/sql"
	"encoding/json"
	_ "log"
	"net/http"
	"todo/database"
	"todo/logging"
)

type Task struct {
	Id   int    `json:"Id" db:"id"`
	Desc string `json:"Desc" db:"description"`
}

// Add godoc
// @Summary Add a new task
// @Description Add a new task for the logged-in user
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body Task true "Task to add"
// @Success 200 {object} map[string]interface{} "Task added successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Error adding task"
// @Router /tasks [post]
func Add(w http.ResponseWriter, r *http.Request) {

	//request
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)

	if err != nil || newTask.Desc == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		logging.Log(err, "Invalid request", "warning", 400, r)
		return
	}

	//extracting session
	cookie, _ := r.Cookie("session_id")

	// extract user from db using cookie
	username := ""
	err = database.TODO.Get(&username, "select username from session where session_id=$1", cookie.Value)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		logging.Log(err, "unauthorized", "error", 401, r)
		return
	}

	// query
	query := `SELECT
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
    END `

	// id selection
	err = database.TODO.Get(&newTask.Id, query, username)
	if err != nil {
		http.Error(w, "Error generating ID", http.StatusInternalServerError)
		logging.Log(err, "Error generating ID", "error", 500, r)
		return
	}

	//insertion
	_, err = database.TODO.Exec("INSERT INTO tasks (id,description,username) VALUES ($1, $2, $3)", newTask.Id, newTask.Desc, username)
	if err != nil {
		http.Error(w, "Error inserting task", http.StatusInternalServerError)
		logging.Log(err, "Error inserting task", "error", 500, r)
		return
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Task added successfully!",
		"task":    newTask,
	}
	json.NewEncoder(w).Encode(response)

	logging.Log(err, "Task added successfully!", "info", 200, r)

}

// List godoc
// @Summary List all tasks
// @Description Get all tasks for the logged-in user
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} []Task "Tasks fetched successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Error fetching tasks"
// @Router /tasks [get]
func List(w http.ResponseWriter, r *http.Request) {

	// extracting session id
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		logging.Log(err, "Unauthorized", "warning", 401, r)
		return
	}

	// Define the query
	query := `
        SELECT t.id, t.description 
        FROM tasks t 
        INNER JOIN session s ON t.username = s.username 
        WHERE s.session_id = $1
    `

	// Define a slice to store tasks
	var tasks []Task

	//fetching data
	err = database.TODO.Select(&tasks, query, cookie.Value)
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		logging.Log(err, "Error fetching tasks", "error", 500, r)
		return
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	if len(tasks) == 0 {
		json.NewEncoder(w).Encode(map[string]string{"message": "No Task Found"})
		logging.Log(err, "No Task Found", "info", 200, r)
		return
	}
	json.NewEncoder(w).Encode(tasks)

	logging.Log(err, "Tasks fetched successfully", "info", 200, r)
}

// Update godoc
// @Summary Update a task
// @Description Update the description of an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body Task true "Task to update"
// @Success 200 {object} map[string]interface{} "Task updated successfully"
// @Failure 400 {object} map[string]string "Invalid task ID or description"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Task not found"
// @Failure 500 {object} map[string]string "Error updating task"
// @Router /tasks [put]
func Update(w http.ResponseWriter, r *http.Request) {

	//extracting id from body
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	id := newTask.Id
	if err != nil || id <= 0 || newTask.Desc == "" {
		http.Error(w, "Invalid task ID or description", http.StatusBadRequest)
		logging.Log(err, "Invalid task ID or description", "warning", 400, r)
		return
	}

	//extracting session
	cookie, _ := r.Cookie("session_id")

	// extract user from db using cookie
	username := ""
	err = database.TODO.Get(&username, "select username from session where session_id=$1", cookie.Value)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		logging.Log(err, "unauthorized", "error", 401, r)
		return
	}

	//updating the task
	var result sql.Result
	result, err = database.TODO.Exec("UPDATE tasks SET description = $2 WHERE id = $1 and username = $3", newTask.Id, newTask.Desc, username)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		logging.Log(err, "Error updating task", "error", 500, r)
		return
	}

	//get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error getting rows affected", http.StatusInternalServerError)
		logging.Log(err, "Error getting rows affected", "error", 500, r)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		logging.Log(err, "Task not found", "warning", 404, r)
		return
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Task updated successfully!",
		"task":    newTask,
	})

	logging.Log(err, "Task updated successfully!", "info", 200, r)

}

// Delete godoc
// @Summary Delete a task
// @Description Delete a task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body Task true "Task to delete"
// @Success 200 {object} map[string]string "Task deleted successfully"
// @Failure 400 {object} map[string]string "Invalid task ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Task not found"
// @Failure 500 {object} map[string]string "Error deleting task"
// @Router /tasks [delete]
func Delete(w http.ResponseWriter, r *http.Request) {

	//extracting id from body
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	id := newTask.Id

	if err != nil || id <= 0 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		logging.Log(err, "Invalid task ID", "warning", 400, r)
		return
	}

	//extracting session
	cookie, _ := r.Cookie("session_id")

	// extract user from db using cookie
	username := ""
	err = database.TODO.Get(&username, "select username from session where session_id=$1", cookie.Value)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		logging.Log(err, "unauthorized", "error", 401, r)
		return
	}

	//removing the task

	var result sql.Result
	result, err = database.TODO.Exec("DELETE FROM tasks WHERE id = $1 and username = $2", id, username)
	if err != nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		logging.Log(err, "Error deleting task", "error", 500, r)
		return
	}

	//get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error getting rows affected", http.StatusInternalServerError)
		logging.Log(err, "Error getting rows affected", "error", 500, r)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		logging.Log(err, "Task not found", "warning", 404, r)
		return
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})

	logging.Log(err, "Task deleted successfully", "info", 200, r)

}
