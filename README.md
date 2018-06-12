# About the project

A basic Web Analytics API implemented using Go and Redis.

The API lets a site make POST requests to log page hits to Redis. And GET requests by users to get a count of hits logged between two dates or a count of unique visits, and get JSON responses back:

**Log a hit**

* POST: http://localhost:8080/setpageview?clientid=287954214567&guid=L4735E3A265E16EEE03F59718B9B5D03019C07D8B6C51F90DA3A666EEC13AB35&path=/index.html&ref=instagram.com

Get all visits to my site between any two arbitrary start and end times
* GET: http://localhost:8080/getpageviews?s=1528329600&e=1528502399

Get all visits to a particular page on my site between any two arbitrary start and end times 
* GET: http://localhost:8080/getpageviews?s=1528329600&e=1528502399&path=/checkout.html

Get all visits to my site by a particular referer between any two arbitrary start and end times
* GET: http://localhost:8080/getpageviews?s=1528329600&e=1528502399&ref=instagram.com

Get all visits to a particular page on my site by a particular referer between any two arbitrary start and end times
* GET: http://localhost:8080/getpageviews?s=1528329600&e=1528502399&ref=instagram.com&path=/index.html

Get Daily, Monthly, and Yearly unique visitors to my site given a date
* GET: http://localhost:8080/getuniques?d=2018/06/08

## Code Walkthrough

A detailed [screen cast](https://www.youtube.com/watch?v=53Hzt7b2fqc) of getting the app up and running and a detailed code walkthrough:

## Pakages Used
* gomodule/redigo: Go client for Redis (https://github.com/gomodule/redigo)
* gorilla/mux: Go request router (https://github.com/gorilla/mux)
* unrolled/render: Go package for easily rendering JSON, among other formats (https://github.com/unrolled/render)

## Redis concepts used
* Transactions: (https://redis.io/topics/transactions)
*  Hashes, Sorted Sets, and Hyperloglog: (https://redis.io/topics/data-types)

## How to run the App
* Install Go
* Install Redis and have it running locally at the default port 6379
* cd into the project directory and run: ``go run main.go``
* Start making POST & GET requests

## Todos

* Add tests
* Easily able to add more dimensions to the app beyond referer and path
* Change the getpageviews GET call to accept date format of YYYY/MM/DD, instead of UNIX_TIMESTAMP
