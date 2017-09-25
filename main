package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"os"
)

// PayLoad that holds info
type PayLoad struct {
	Project   interface{}
	Owner     interface{}
	Committer interface{}
	Commits   interface{}
	Languages []string
}

//Checks for errors
func err(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	return resp.Body
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	var payLoad PayLoad
	var commit []interface{}

	mResp := new(map[string]interface{})
	languageResp := new(map[string]interface{})

	// Adds full_name and owner.login to payLoad
	json.NewDecoder(err("https://api.github.com/repos" + r.URL.Path)).Decode(mResp)
	payLoad.Project = (*mResp)["full_name"].(string)
	payLoad.Owner = (((*mResp)["owner"].(map[string]interface{}))["login"]).(string)

	// Adds Committer and total commits to payLoad
	json.NewDecoder(err((*mResp)["contributors_url"].(string))).Decode(&commit)
	payLoad.Committer = (commit[0].(map[string]interface{})["login"]).(string)
	payLoad.Commits = commit[0].(map[string]interface{})["contributions"]

	// Adds all languages to payLoad, taking out only keys from map
	json.NewDecoder(err((*mResp)["languages_url"].(string))).Decode(&languageResp)
	for r := range *languageResp {
		payLoad.Languages = append(payLoad.Languages, r)
	}
	// Sort languages alphabetically
	sort.Strings(payLoad.Languages)

	//Print payLoad
	m, _ := json.MarshalIndent(payLoad, "", "    ")
	fmt.Fprint(w, string(m))

}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080
	}
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
