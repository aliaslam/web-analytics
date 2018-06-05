package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	. "github.com/aliaslam/webanalytics/utils"
)

/*
Sample POST request: http://localhost:8080/?clientid=28795421456&guid=1041e0c3a-6c58-4160-9196-f0050faf00ef&path=/checkout.html&ref=google.com
Response: 200 OK
*/

// AddPageview accepts a POST request and presist the hit in Redis
func SetPageview(writer http.ResponseWriter, request *http.Request) {

	if request.Method == "POST" {
		request.ParseForm()

		//Gather all the POST params
		clientid := request.FormValue("clientid")
		guid := request.FormValue("guid")
		path := request.FormValue("path")
		ref := request.FormValue("ref")

		t := time.Now()
		unixTimestamp := t.Unix()

		eventType := "pageview" //Currently we are logging the pageview event but we can extend the system to log any type of event
		pageviewIndexKey := clientid + KS + "pageviewindex"
		pageViewIndex, err := RC.Do("INCR", pageviewIndexKey) //Sets up the auto-incr index
		if err != nil {
			log.Println(err)
		}

		//Setup the format of all the keys
		pageviewHashKey := clientid + KS + eventType + KS + fmt.Sprintf("%v", pageViewIndex)
		pageviewByPathKey := clientid + KS + "path" + KS + path
		pageviewByRefKey := clientid + KS + "ref" + KS + ref
		uniquesKey := clientid + KS + "uniques" + KS + t.Format("2006:01:02")
		timeIndexKey := clientid + KS + "timeindex"

		//Create/Update the hash, sorted sets, and hyperloglog for uniques in a transaction
		RC.Send("MULTI")
		RC.Send("HMSET", pageviewHashKey, "guid", guid, "path", path, "ref", ref)
		RC.Send("ZADD", timeIndexKey, unixTimestamp, pageViewIndex)
		RC.Send("ZADD", pageviewByPathKey, unixTimestamp, pageViewIndex)
		RC.Send("ZADD", pageviewByRefKey, unixTimestamp, pageViewIndex)
		RC.Send("PFADD", uniquesKey, guid)
		_, err = RC.Do("EXEC")

		if err != nil {
			log.Println(err)
		}
	}

}
