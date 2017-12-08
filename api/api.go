package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"sort"
	"time"

)


func (p pointSlice) Len() int {
	return len(p)
}

func (p pointSlice) Less(i, j int) bool {
	return p[i].Date.Before(p[j].Date)
}

func (p pointSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func reverse(ss []Issue) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}

func (api *API) GenerateIssueChart(repoString string) (IssueChart){
		var a IssueChart
		var open Dataset
		var closed Dataset
		res := api.Database.Find("issue/" + repoString)
		if(res != ""){
			byteRes := []byte(res)
			err := json.Unmarshal(byteRes, &a)
			if err != nil {
				fmt.Printf("%v", err.Error())
			}
		}else{
			open.Label = "Open Issues"
			closed.Label = "Closed Issues"

			a.Name = repoString
			repo := api.GetRepo(repoString)
			startTime := time.Now()

			for _, issue := range repo.Issues {
				var openTime time.Duration
				var point Point
				ignore := false
				for _,label := range issue.Labels{
					if(label.Name == "issue/ignore"){
						ignore = true
					}
				}

				if(ignore){continue}

				point.Link = issue.URL
				if issue.State == "open" {
					a.Open += 1
					openTime = startTime.Sub(issue.Created) / time.Second
					point.Label = issue.Name + " - " + strconv.Itoa(issue.Number)
					point.Value = int64(openTime)
					open.Points = append([]Point{point}, open.Points...)
					} else {
						openTime = issue.Closed.Sub(issue.Created) / time.Second
						point.Label = issue.Name + " - " + strconv.Itoa(issue.Number)
						point.Value = int64(openTime)
						closed.Points = append([]Point{point}, closed.Points...)
						a.Closed += 1
					}

					a.AvgDuration += openTime
					if openTime > a.MaxDuration {
						a.MaxDuration = openTime
					}

				}

				a.Data = append(a.Data, open)
				a.Data = append(a.Data, closed)

				a.AvgDuration /= time.Duration(a.Open + a.Closed)
				api.Database.Set("issue/" + repoString, a);
				api.Database.Expire("issue/" + repoString, api.Database.TTL(repoString)/time.Second)
		}
		return a;
}

func (api *API) GenerateStaleness(repoString string) (Staleness){
	a := api.GenerateIssueChart(repoString)
	var stl Staleness

	for _,issue := range a.Data[0].Points{
		stl.Stale += issue.Value;
	}
	if(len(a.Data[0].Points) > 0){
		stl.Stale /= int64(len(a.Data[0].Points));
	}
	stl.Max = int64(a.MaxDuration)
	stl.Ratio = float32(stl.Stale)/float32(stl.Max)
	if(stl.Ratio >= .75){
		stl.Text = "Looking pretty stale"
	}else if(stl.Ratio >= .5){
		stl.Text = "Slightly stale"
	}else{
		stl.Text = "Looking good"
	}
	return stl
}

func (api *API) GetStaleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoString := vars["owner"] + "/" + vars["repo"]

	stl := api.GenerateStaleness(repoString)
	WriteJSON(w, stl)

}

func (api *API) GetBadgeHandler(w http.ResponseWriter, r *http.Request) {
 	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")

	vars := mux.Vars(r)
	repoString := vars["owner"] + "/" + vars["repo"]
	stl := api.GenerateStaleness(repoString)

	if(stl.Ratio >= .75){
		http.ServeFile(w, r, "stale.svg")
	}else if(stl.Ratio >= .5){
		http.ServeFile(w, r, "getting_stale.svg")
	}else{
		http.ServeFile(w, r, "looking_good.svg")
	}
}

func (api *API) GetIssueChart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoString := vars["owner"] + "/" + vars["repo"]
	a := api.GenerateIssueChart(repoString)
	WriteJSON(w, a)
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

func (api *API) GetBarChart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoString := vars["owner"] + "/" + vars["repo"]
	repo := api.GetRepo(repoString)
	var chart pointSlice
	for _, issue := range repo.Issues {
		str := ""
		weight := false
		for _,label := range issue.Labels{
			if(strings.Contains(label.Name, "burndown")){
				weight = true
				str = label.Name[9:len(label.Name)]
				break;
			}
		}
		var val int
		val = 1
		if(weight){
			val,_ = strconv.Atoi(str)
		}
		chart = append(chart,Point{Label:strconv.Itoa(issue.Number),Value:int64(val),Date:issue.Created})
		if(issue.State =="closed"){
			chart = append(chart,Point{Label:strconv.Itoa(issue.Number),Value:int64(-val),Date:issue.Closed})

		}
	}
	sort.Sort(chart);
	WriteJSON(w, chart)
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
		if err != nil {
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
		api.Database.Expire(data, 6000)
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
