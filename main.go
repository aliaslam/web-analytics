package main

import (
	"log"
	"net/http"

	"github.com/aliaslam/webanalytics/controllers"
	"github.com/aliaslam/webanalytics/utils"
	"github.com/gorilla/mux"
)

//main() func that sets up the redis connection, sets up the REST API routes, and starts the HTTP server
func main() {
	utils.RC = utils.GetRedisConnection()
	defer utils.RC.Close()
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.SetPageview).Methods("POST")
	r.HandleFunc("/pageviews", controllers.GetPageview).Methods("GET")
	r.HandleFunc("/uniques", controllers.GetUnique).Methods("GET")

	log.Println("(Web Analytics v1.0) : Listening on localhost:8080 for requests...")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
