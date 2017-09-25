package main

import (
"net/http"
"log"
"encoding/json"
"fmt"
)

type RepoInfo struct {
	Project 	string 		`json:"full_name"`
	Owner		OwnerStruct 	`json:"owner`
	LanguagesUrl 	string		`json:"languages_url"`
	ContributorsUrl string 		`json:"contributors_url"`
	TopContributor	Contributor
	Languages 	[]interface{}

}

type OwnerStruct struct {
	Login string `json:"login"`
}

type Contributor struct {
	Login string `json:"login"`
	Contributions int `json:"contributions"`
}

func getRepoInfo(w *http.ResponseWriter, r *http.Request, dest *RepoInfo) {
	resp, err := http.Get("https://api.github.com/repos" + r.URL.Path)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		json.NewDecoder(resp.Body).Decode(&dest)
	}
}

func getLanguages(w *http.ResponseWriter, r *http.Request, languageUrl string, dest *RepoInfo) {
	resp, err := http.Get(languageUrl)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		langs := new(map[string]interface{})
		json.NewDecoder(resp.Body).Decode(&langs)

		for r := range *langs {
			dest.Languages = append(dest.Languages, r)
		}
	}
}

func getTopContributor(w *http.ResponseWriter, r *http.Request, dest *RepoInfo) {
	resp, err := http.Get(dest.ContributorsUrl)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		var contributors []Contributor
		json.NewDecoder(resp.Body).Decode(&contributors)

		var topContributor Contributor
		for i := 0; i < len(contributors); i++ {
			if contributors[i].Contributions > topContributor.Contributions {
				topContributor = contributors[i]
			}
		}
		dest.TopContributor.Login = topContributor.Login
		dest.TopContributor.Contributions = topContributor.Contributions
	}
}


func handlerFunc(w http.ResponseWriter, r *http.Request) {
	var repoInfo RepoInfo
	getRepoInfo(&w, r, &repoInfo)
	getLanguages(&w, r, repoInfo.LanguagesUrl, &repoInfo)
	getTopContributor(&w, r, &repoInfo)
	m, _ := json.MarshalIndent(repoInfo, "", "    ")
	fmt.Fprint(w, string(m))
	fmt.Fprintln(w,"hello")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe("127.0.0.1:8083", nil)
}
