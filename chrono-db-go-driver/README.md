# Sample Citra db Go Driver Usage

This repo is a demo on how to use chrono-db-go-driver


### Pre-Requisite

1. Keep chrono db local setup ready
2. Create the chrono db with admin access & save the creds
3. Start the chrono db with admin access

### Steps to use chrono db

1. Copy the main.go file
2. Pass the Keeper & secret you get after creating chrono-db
3. go mod tidy or run
```
go get github.com/citra-org/chrono-db-go-driver
```
4. Run
```
go run main.go
```
5. After the Server starts
check the api.curl for the curl commands

 To create a stream
 ```
 curl http://localhost:3000/cs/<stream> 
 ```
 To write a event into a stream

 ```
    curl -X POST -d '{"<header>": "<body>","<header>": "<body>"}' -H "Content-Type: application/json" http://localhost:3000/w/<stream>
 ```
 To read the events from a stream

 ```
 curl http://localhost:3000/r/<stream>
 ```


**NOTE: THIS CODEBASE IS STILL UNDER DEV & HAS LOT OF BUGS & IMPROVEMENTS TO BE DONE, PLEASE DONT USE IN PROD**