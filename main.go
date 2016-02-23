package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

type ClassInfo struct {
	Credits int
	Prereqs []string
	Title   string
	Course  string
}

var (
	ErrNoSuchClass = fmt.Errorf("No Such Class")
)

var pageCache map[string]string

func init() {
	pageCache = make(map[string]string)
}

func GetClassPage(class string) (string, error) {

	c := http.Client{}

	if _, isOk := pageCache[class]; isOk {
		return pageCache[class], nil
	}

	resp, err := c.Get("http://www.registrar.ufl.edu/cdesc?crs=" + class)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(resp.Body)

	pageCache[class] = string(data)

	return string(data), nil
}

func ParsePrereqs(class string, prereqs string) []string {
	parsed := make([]string, 0)

	re := regexp.MustCompile("(?P<class>[A-Z]{3} [0-9]{4})")

	allPrereqs := re.FindAllString(prereqs, -1)

	parsed = append(allPrereqs[:])

	final := make([]string, 0)

	for _, item := range parsed {
		item = item[:3] + item[4:]

		if item != class {
			final = append(final, item)
		}
	}

	return final
}

func ParseClassPage(data string) (ci ClassInfo, err error) {
	defer func() {
		// there was no such class page
		if r := recover(); r != nil {
			err = nil
		}
	}()

	classRE := regexp.MustCompile("<h2>(?P<class>.*)</h2>")
	courseRE := regexp.MustCompile("<h3>(?P<title>.*)</h3>")
	creditsRE := regexp.MustCompile("<strong>Credits: (?P<credits>\\d*)")
	prereqRE := regexp.MustCompile("Prereq: (?P<prereq>.*)</strong>")

	ci.Title = classRE.FindStringSubmatch(data)[1]
	ci.Course = courseRE.FindStringSubmatch(data)[1]
	creditsInt := creditsRE.FindStringSubmatch(data)
	prereqInt := prereqRE.FindStringSubmatch(data)

	ci.Credits, _ = strconv.Atoi(creditsInt[1])
	ci.Course = ci.Course[0:3] + ci.Course[4:len(ci.Course)-1]
	ci.Prereqs = ParsePrereqs(ci.Course, prereqInt[1])

	return ci, nil
}

func GetClassInfo(class string) (ClassInfo, error) {
	data, err := GetClassPage(class)
	if err != nil {
		fmt.Println("error: invalid class page", class)
		return ClassInfo{}, err
	}

	ci, err := ParseClassPage(data)
	if err != nil {
		fmt.Println("error: invalid parsing of class page", class)
		return ClassInfo{}, err
	}

	return ci, nil
}

func GetClassInfoJSON(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	class := vars["class"]
	data, _ := GetClassPage(class)

	ci, _ := ParseClassPage(data)

	json.NewEncoder(w).Encode(ci)
}

type TemplateData struct {
	Met   string
	Unmet string
}

func StartCalc(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting Calculation")

	tmpl, err := template.ParseFiles("templates/results.html")
	if err != nil {
		http.Error(w, "template syntax incorrect: "+err.Error(), 500)
		return
	}

	r.ParseForm()

	met := r.FormValue("met")
	unmet := r.FormValue("unmet")

	tmpl.Execute(w, &TemplateData{Met: met, Unmet: unmet})
}

func ServeFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving File: ", path.Clean(r.URL.RequestURI()))

	uri := "./templates/" + path.Clean(r.URL.RequestURI())
	http.ServeFile(w, r, uri)
}

type statistics struct {
	UsersHelped    int `json:"users_helped"`
	AvgTimeTaken   int `json:"avg_time_taken_ms"`
	TotalTimeTaken int `json:"total_time_taken_ms"`
}

func GenerateStats(w http.ResponseWriter, r *http.Request) {
	log.Println("Generating Stats")

	var stats statistics

	stats.TotalTimeTaken = TotalTimeTaken
	stats.AvgTimeTaken = TotalTimeTaken / UsersHelped
	stats.UsersHelped = UsersHelped

	json.NewEncoder(w).Encode(stats)
}

func main() {

	address := ":8888"

	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/class/{class}", GetClassInfoJSON)
	router.HandleFunc("/api/startCalc", StartCalc)
	router.HandleFunc("/api/calc", CalcClassesJSON)
	router.HandleFunc("/api/stats", GenerateStats)
	router.PathPrefix("/").HandlerFunc(ServeFile)

	log.Fatal(http.ListenAndServe(address, router))
}
