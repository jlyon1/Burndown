package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"fmt"
	"burndown/database"
	"github.com/bradfitz/slice"
)

type API struct {
	Database database.DB
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
	Label time.Time
	Value int
}


func reverse(ss []Issue) {
    last := len(ss) - 1
    for i := 0; i < len(ss)/2; i++ {
        ss[i], ss[last-i] = ss[last-i], ss[i]
    }
}

func (api *API) GetRepo(data string) []Point{
	url := "https://api.github.com/repos/" + data + "/issues?state=all"
	_ = url
	var labels []Point
	res := api.Database.Find(data);
	if res != ""{
		byteRes := []byte(res)
		err:= json.Unmarshal(byteRes,&labels)
		if err != nil{
			fmt.Printf("%v",err.Error())
		}
		fmt.Printf("\nCacheload\n");
	}else{
		var issues []Issue
		resp, err := http.Get(url)
		if err != nil{
			fmt.Printf("%v",err.Error())

		}
		reader := json.NewDecoder(resp.Body)
		reader.Decode(&issues)
		for idx, _ := range issues{
			issues[idx].Weight = 1;
		}

		//var labels []Point

		for _, issue := range issues {
			labels = append(labels, Point{Label: issue.Created,Value: issue.Weight})

			if(issue.State == "closed"){
				labels = append(labels, Point{Label: issue.Closed,Value: - issue.Weight})

			}
		}
		slice.Sort(labels, func(i,j int) bool {
			return labels[i].Label.After(labels[j].Label)
		})
		api.Database.Set(data,labels)
		api.Database.Expire(data,100);
	}

	return labels;
}

func (api *API) GetRepoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoString := vars["owner"] + "/" + vars["repo"]
	labels := api.GetRepo(repoString)
	fmt.Printf("ASDF %d\n",len(labels));
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
