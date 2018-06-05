package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	. "github.com/aliaslam/webanalytics/utils"
	"github.com/gomodule/redigo/redis"
	"github.com/unrolled/render"
)

/*
Sample GET requests and JSON responses:
http://localhost:8080/pageviews?guid=862e0c3a-6c58-4160-9196-f0050faf00ef&path=/checkout.html&ref=instagram.com&s=1525132800&e=1527724800
{"Views":13,"StartTime":"1525132800","EndTime":"1527724800","Path":"/checkout.html","Ref":"instagram.com"}

http://localhost:8080/pageviews?guid=862e0c3a-6c58-4160-9196-f0050faf00ef&ref=google.com&s=1525132800&e=1527724800
{"Views":93,"StartTime":"1525132800","EndTime":"1527724800","Path":"ALL","Ref":"google.com"}

http://localhost:8080/pageviews?guid=862e0c3a-6c58-4160-9196-f0050faf00ef&path=/index.html&s=1525132800&e=1527724800
{"Views":55,"StartTime":"1525132800","EndTime":"1527724800","Path":"/index.html","Ref":"ALL"}

http://localhost:8080/pageviews?guid=862e0c3a-6c58-4160-9196-f0050faf00ef&s=1527379200&e=1527465599
{"Views":200,"StartTime":"1527379200","EndTime":"1527465599","Path":"ALL","Ref":"ALL"}

http://localhost:8080/uniques?t=1527490930
{"Daily":5,"Monthly":5,"Yearly":5,"Time":"2018:05:28"}
*/

type pageviews struct {
	Views     int64
	StartTime string
	EndTime   string
	Path      string
	Ref       string
}

type uniques struct {
	Daily   int64
	Monthly int64
	Yearly  int64
	Time    string
}

var clientid = "28795421456" //This would be fetched from the DB or the client's session
var pv *pageviews
var r = render.New()

//Returns the Daily, Monthly, and Yearly Unique hits
func GetUnique(writer http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		t := request.URL.Query().Get("t")

		parsedUT, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			log.Println(err)
		}
		tm := time.Unix(parsedUT, 0)
		date := tm.Format("2006:01:02")

		dailyUniques, err := RC.Do("PFCOUNT", clientid+KS+"uniques"+KS+date) //PFCOUNT on the hyperloglog object for daily uniques
		if err != nil {
			log.Println(err)
		}

		monthlyUniques := queryUniques("month", tm)
		yearlyUniques := queryUniques("year", tm)

		uniques := uniques{
			Daily:   dailyUniques.(int64),
			Monthly: monthlyUniques.(int64),
			Yearly:  yearlyUniques.(int64),
			Time:    date,
		}
		r.JSON(writer, http.StatusOK, uniques)

	}
}

//Helper function to get the Monthly, and Yearly Unique hits
func queryUniques(interval string, tm time.Time) interface{} {

	var timeformat, key string

	switch interval {
	case "month":
		timeformat = "2006:01"
		key = "monthlyuniques"
	case "year":
		timeformat = "2006"
		key = "yearlyuniques"
	}
	//Depending on monthly or yearly query, we first PFMERGE all the HLL keys, then do a PFCOUNT to get the uniques
	matchingKeys := getMatchingKeys(clientid + KS + "uniques" + KS + tm.Format(timeformat) + KS + "*")

	combniedKeys := append([]string{clientid + KS + key}, matchingKeys...)

	s := make([]interface{}, len(combniedKeys))
	for index, value := range combniedKeys {
		s[index] = value
	}

	_, err := RC.Do("PFMERGE", s...)
	if err != nil {
		log.Println(err)
	}

	uniques, err := RC.Do("PFCOUNT", clientid+KS+key)
	if err != nil {
		log.Println(err)
	}

	return uniques
}

//Helper function to get mathcing keys from redis based on the given pattern
func getMatchingKeys(pattern string) []string {
	iter := 0

	keys := []string{}
	for {

		if arr, err := redis.Values(RC.Do("SCAN", 0, "MATCH", pattern, "COUNT", 365)); err != nil {
			panic(err)
		} else {

			iter, _ = redis.Int(arr[0], nil)
			keys, _ = redis.Strings(arr[1], nil)
		}
		if iter == 0 {
			break
		}
	}
	return keys
}

// GetPageview responds to a GET request with page view info
func GetPageview(writer http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {

		getParams := pageviews{
			Views:     0,
			StartTime: request.URL.Query().Get("s"),
			EndTime:   request.URL.Query().Get("e"),
			Path:      request.URL.Query().Get("path"),
			Ref:       request.URL.Query().Get("ref"),
		}

		if getParams.Path != "" || getParams.Ref != "" {

			if getParams.Path != "" && getParams.Ref != "" {
				pv = getPageviewsByRefAndPath(&getParams)
			}

			if getParams.Path == "" && getParams.Ref != "" {
				pv = getPageviewsByRef(&getParams)
			}

			if getParams.Path != "" && getParams.Ref == "" {
				pv = getPageviewsByPath(&getParams)
			}

		} else {
			pv = getAllPageViews(&getParams)
		}

		r.JSON(writer, http.StatusOK, pv)
	}

}

//Does an intersection of the views and paths sorted sets to get page views for a given ref and path
func getPageviewsByRefAndPath(pv *pageviews) *pageviews {
	hits, err := RC.Do("ZINTERSTORE", "out", 2, clientid+KS+"path"+KS+pv.Path, clientid+KS+"ref"+KS+pv.Ref, "WEIGHTS", pv.StartTime, pv.EndTime)
	if err != nil {
		log.Println(err)
	}
	pv.Views = hits.(int64)
	return pv
}

//Does a ZCOUNT on the ref sorted set to get page views for a given ref
func getPageviewsByRef(pv *pageviews) *pageviews {
	hits, err := RC.Do("ZCOUNT", clientid+KS+"ref"+KS+pv.Ref, pv.StartTime, pv.EndTime)
	if err != nil {
		log.Println(err)
	}
	pv.Views = hits.(int64)
	pv.Path = "ALL"
	return pv
}

//Does a ZCOUNT on the path sorted set to get page view for a given path
func getPageviewsByPath(pv *pageviews) *pageviews {
	hits, err := RC.Do("ZCOUNT", clientid+KS+"path"+KS+pv.Path, pv.StartTime, pv.EndTime)
	if err != nil {
		log.Println(err)
	}
	pv.Views = hits.(int64)
	pv.Ref = "ALL"
	return pv
}

//Does a ZCOUNT on the timeindex sorted set to get all views
func getAllPageViews(pv *pageviews) *pageviews {
	hits, err := RC.Do("ZCOUNT", clientid+KS+"timeindex", pv.StartTime, pv.EndTime)
	if err != nil {
		log.Println(err)
	}
	pv.Views = hits.(int64)
	pv.Path = "ALL"
	pv.Ref = "ALL"
	return pv
}
