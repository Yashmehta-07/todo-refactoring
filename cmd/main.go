package main

import (
	"fmt"
	"log"
	"net/http"
	_ "strconv"
	_ "time"
	"todo/database"
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
	fmt.Println("Server running on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))

}
