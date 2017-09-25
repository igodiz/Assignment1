package main

import (
	"net/http"
	"log"
	"encoding/json"
	"fmt"
)

type RepoInfo struct {
	Project 	string 		`json:"full_name"`
	Owner		OwnerStruct 	`json:"owner"`
	LanguagesUrl 	string		`json:"languages_url"`
	ContributorsUrl string 		`json:"contributors_url"`
	Languages 	[]interface{}
	TopContributor	Contributor
}

type OwnerStruct struct {
	Login string `son:"login"`
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

func PrettyPrint(w http.ResponseWriter, info RepoInfo) {
	fmt.Fprint(w, "{\n")
	fmt.Fprint(w, "\t\"project\": \"" + info.Project + "\",\n")
	fmt.Fprint(w, "\t\"owner\": \"" + info.Owner.Login + "\",\n")
	fmt.Fprint(w, "\t\"committer\": \"" + info.TopContributor.Login + "\",\n")
	fmt.Fprintf(w,"\t\"commits\": \"%d\"%s"  ,info.TopContributor.Contributions, ",\n")
	fmt.Fprintf(w, "\t\"languages\": \"%v%s"  ,info.Languages, "\"\n")
	fmt.Fprint(w, "}\n")
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	var info RepoInfo
	getRepoInfo(&w, r, &info)
	getLanguages(&w, r, info.LanguagesUrl, &info)
	getTopContributor(&w, r, &info)
	PrettyPrint(w, info)
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
