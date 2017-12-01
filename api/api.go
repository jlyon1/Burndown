package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type API struct {
}

type Repo struct {
	Name      string `json:"name"`
	FullName  string `json:"full_name"`
	URL       string `json:"html_url"`
	IssuesURL string `json:"issues_url"`
	PullsURL  string `json:"pulls_url"`
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
	Weight  int
}

type Repository struct {
	Name   string
	Issues []Issue
	URL    string
}

type Point struct{
	Labels []int
	Value []int
}

func reverse(ss []Issue) {
    last := len(ss) - 1
    for i := 0; i < len(ss)/2; i++ {
        ss[i], ss[last-i] = ss[last-i], ss[i]
    }
}

func (api *API) GetRepoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := "https://api.github.com/repos/" + vars["repo"] + "/" + vars["owner"] + "/issues?state=all"
	var issues []Issue
	resp, _ := http.Get(url)
	reader := json.NewDecoder(resp.Body)
	reader.Decode(&issues)
	reverse(issues)
	for idx, _ := range issues{
		issues[idx].Weight = 1;
		if(issues[idx].State == "closed" && idx != 0){
			issues[idx].Weight = -issues[idx].Weight
		}
	}
	var labels []int
	var values []int
	for idx, issue := range issues {
		labels = append(labels, issue.Number)
		if(idx > 0){
			values = append(values, issue.Weight + values[idx - 1])
		}else{
			values = append(values, issue.Weight)
		}
	}
	data := Point{
		Labels: labels,
		Value: values,
	}
	WriteJSON(w, data)
}

func (api *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func WriteJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Write(b)
	return nil
}
