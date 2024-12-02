package main

import (
	"net/http"
	_ "strconv"
	_ "time"
	"todo/database"
	"todo/logging"
	"todo/routes"

	_ "todo/docs"

	_ "github.com/lib/pq" // Import pq driver
)

// @title           To-Do API
// @version         1.0
// @description     A brief description of your API
// @host            localhost:8000
// @BasePath        /
func main() {

	//initializing database
	db := database.ConnectDB()
	//closing database onces the server is closed
	defer db.Close()

	// // share the DB to auth package
	// utils.SetDB(db)

	r := routes.Route()

	//server start
	logging.Log(nil, "Server running on http://localhost:8000", "info", 200, nil)
	logging.Logger.Fatal(http.ListenAndServe(":8000", r))

}
