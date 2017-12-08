package main

import (
	"burndown/api"
	"github.com/gorilla/mux"
	"net/http"
	"burndown/database"
	"fmt"
)

func connectDB(db database.DB) {
	fmt.Printf("Connected: %v",db.Connect())

}

func main() {
	redis := &database.Redis{}
	redis.IP = "localhost"
	redis.Port = "6379"
	redis.DB = 0
	redis.Password = ""
	connectDB(redis)

	api := api.API{
		Database: redis,
	}

	r := mux.NewRouter()

	r.HandleFunc("/", api.IndexHandler).Methods("GET")
	r.HandleFunc("/{owner}/{repo}", api.IndexHandler).Methods("GET")

	r.HandleFunc("/get/{owner}/{repo}",api.GetRepoHandler).Methods("GET")
	r.HandleFunc("/stale/{owner}/{repo}",api.GetStaleHandler).Methods("GET")
	r.HandleFunc("/badge/{owner}/{repo}",api.GetBadgeHandler).Methods("GET")

	r.HandleFunc("/issue/{owner}/{repo}",api.GetIssueChart).Methods("GET")
	r.HandleFunc("/bar/{owner}/{repo}",api.GetBarChart).Methods("GET")
	r.HandleFunc("/valid/{owner}/{repo}",api.ValidHandler).Methods("GET")

	r.PathPrefix("/static/serve/").Handler(http.StripPrefix("/static/serve/", http.FileServer(http.Dir("static/"))))
	http.ListenAndServe("0.0.0.0:8080", r)

}
