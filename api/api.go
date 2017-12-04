package api

import (
	"burndown/database"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"strconv"
	"strings"
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

type Point struct {
	Label string
	Value int64
}

type Dataset struct {
	Label  string
	Points []Point
}

type Chart struct {
	Data []Dataset
}

type IssueChart struct {
	Name 				string
	Data        []Dataset
	Open        int
	Closed      int
	AvgDuration time.Duration
	MaxDuration time.Duration
}

func reverse(ss []Issue) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}

func (api *API) GetIssueChart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoString := vars["owner"] + "/" + vars["repo"]

	var a IssueChart
	var open Dataset
	var closed Dataset

	open.Label = "Open Issues"
	closed.Label = "Closed Issues"


	a.Name = repoString
	repo := api.GetRepo(repoString)
	startTime := time.Now()

	for _, issue := range repo.Issues{
		var openTime time.Duration
		var point Point

		if(issue.State == "open"){
			a.Open += 1;
			openTime = startTime.Sub(issue.Created)/time.Second
			point.Label = issue.Name + " - " + strconv.Itoa(issue.Number)
			point.Value = int64(openTime)
			open.Points = append([]Point{point},open.Points...)
		}else{
			openTime = issue.Closed.Sub(issue.Created)/time.Second
			point.Label = issue.Name + " - " + strconv.Itoa(issue.Number)
			point.Value = int64(openTime)
			closed.Points = append([]Point{point},closed.Points...)
			a.Closed += 1;
		}

		a.AvgDuration += openTime
		if(openTime > a.MaxDuration){
			a.MaxDuration = openTime
		}

	}

	a.Data = append(a.Data,open)
	a.Data = append(a.Data,closed)

	a.AvgDuration /= time.Duration(a.Open + a.Closed)
	WriteJSON(w,a)
}

func (api *API) ValidHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	repoString := vars["owner"] + "/" + vars["repo"]
	res := api.Database.Find(repoString)
	url := "https://api.github.com/repos/" + repoString
	if res != "" {
		WriteJSON(w, "true")
	} else {
		resp, err := http.Get(url)
		var repo Repository
		if err != nil {
			fmt.Printf("%v", err.Error())
		}
		reader := json.NewDecoder(resp.Body)
		reader.Decode(&repo)
		if repo.Name != "" {
			WriteJSON(w, "true")
			api.GetRepo(repoString)
		} else {
			WriteJSON(w, "false")
		}
	}

}

func (api *API) GetRepo(data string) Repository {
	data = strings.ToLower(data)
	url := "https://api.github.com/repos/" + data

	_ = url
	var repo Repository
	res := api.Database.Find(data)
	if res != "" {
		byteRes := []byte(res)
		err := json.Unmarshal(byteRes, &repo)
		if err != nil {
			fmt.Printf("%v", err.Error())
		}
	} else {

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("%v", err.Error())
		}
		reader := json.NewDecoder(resp.Body)
		reader.Decode(&repo)

		issue := url + "/issues?state=all&per_page=100"
		resp, err = http.Get(issue)
		if err != nil {
			fmt.Printf("%v", err.Error())
		}
		reader = json.NewDecoder(resp.Body)
		reader.Decode(&repo.Issues)

		commits := url + "/commits?state=all&per_page=100"
		resp, err = http.Get(commits)
		if err != nil{
			fmt.Printf("%v", err.Error())
		}
		reader = json.NewDecoder(resp.Body)
		reader.Decode(&repo.Commits)

		Pulls := url + "/pulls?state=all&per_page=100"
		resp, err = http.Get(Pulls)
		if err != nil {
			fmt.Printf("%v", err.Error())
		}
		reader = json.NewDecoder(resp.Body)
		reader.Decode(&repo.Pulls)

		api.Database.Set(data, repo)
		api.Database.Expire(data, 10000)
	}

	return repo
}

func (api *API) GetRepoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoString := vars["owner"] + "/" + vars["repo"]
	labels := api.GetRepo(repoString)
	WriteJSON(w, labels)
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
