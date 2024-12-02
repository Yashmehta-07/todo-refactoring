package main

import (
	"net/http"
	_ "strconv"
	_ "time"
	"todo/database"
	"todo/logging"
	"todo/routes"

	_ "github.com/lib/pq" // Import pq driver
)

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
