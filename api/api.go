package api

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

type API struct {

}

type Repo struct{
	Name string `json:"name"`
	FullName string `json:"full_name"`
	URL string `json:"html_url"`
	IssuesURL string `json:"issues_url"`
	PullsURL string `json:"pulls_url"`
}

func (api *API) GetRepoHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	url := "https://api.github.com/repos/" + vars["repo"] + "/" + vars["owner"];
	var test Repo
	resp,err := http.Get(url)
	_ = err;
 	reader := json.NewDecoder(resp.Body)
	reader.Decode(&test)
	// defer resp.Body.Close();
	// body, _ := ioutil.ReadAll(resp.Body);
	WriteJSON(w,test);
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
