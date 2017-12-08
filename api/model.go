package api

import (
  "burndown/database"
  "time"

)

type API struct {
	Database database.DB
}

type Label struct {
	Name string `json:"name"`
}

type Issue struct {
	Name    string    `json:"title"`
	Number  int       `json:"number"`
	State   string    `json:"state"`
	Created time.Time `json:"created_at"`
	Closed  time.Time `json:"closed_at"`
	Labels  []Label   `json:"labels"`
	URL			string		`json:"html_url"`
	Weight  int
}

type CommitInfo struct {
	Message string `json:"message"`
	URL     string `json:"html_url"`
}

type User struct {
	Name   string    `json:"login"`
	Id     int       `json:"id"`
	Avatar string    `json:"avatar_url"`
	Url    string    `json:"html_url"`
	Date   time.Time `json:"date"`
}

type Commit struct {
	Info      CommitInfo `json:"commit"`
	Author    User       `json:"Author"`
	Committer User       `json:"committer"`
}

type Pull struct {
}

type Repository struct {
	Name    string `json:"full_name"`
	Owner   User   `json:"owner"`
	URL     string `json:"html_url"`
	Issues  []Issue
	Pulls   []Pull
	Commits []Commit
}

type Staleness struct {
	Stale int64 `json:"staleness"`
	Max   int64 `json:"max"`
	Ratio float32	`json:ratio`
	Text 	string	`json:text`
}

type Point struct {
	Label string
	Value int64
	Link   string
	Date  time.Time
}

type Dataset struct {
	Label  string
	Points []Point
}

type Chart struct {
	Data []Dataset
}

type IssueChart struct {
	Name        string
	Data        []Dataset
	Open        int
	Closed      int
	AvgDuration time.Duration
	MaxDuration time.Duration
}

type pointSlice []Point
