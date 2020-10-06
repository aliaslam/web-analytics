package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nathanusask/web-analytics/controllers"
	"github.com/nathanusask/web-analytics/utils"
)

//main() func that sets up the redis connection, sets up the REST API routes, and starts the HTTP server
func main() {
	utils.RC = utils.GetRedisConnection()
	defer utils.RC.Close()
	r := mux.NewRouter()
	r.HandleFunc("/setpageview", controllers.SetPageview).Methods("POST")
	r.HandleFunc("/getpageviews", controllers.GetPageviews).Methods("GET")
	r.HandleFunc("/getuniques", controllers.GetUniques).Methods("GET")

	log.Println("(Web Analytics v1.0) : Listening on localhost:8080 for requests...")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
