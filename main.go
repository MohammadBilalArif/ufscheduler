package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func GetClassPage(class string) (string, error) {

	c := http.Client{}

	resp, err := c.Get("http://www.registrar.ufl.edu/cdesc?crs=" + class)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(resp.Body)

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

		fmt.Printf("Item: '%s'\n", item)

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
			err = ErrNoSuchClass
		}
	}()

	classRE := regexp.MustCompile("<h2>(?P<class>.*)</h2>")
	courseRE := regexp.MustCompile("<h3>(?P<title>.*)</h3>")
	prereqRE := regexp.MustCompile("<strong>Credits: (?P<credits>\\d*); Prereq: (?P<prereq>.*)</strong>")

	ci.Title = classRE.FindStringSubmatch(data)[1]
	ci.Course = courseRE.FindStringSubmatch(data)[1]
	prereqInt := prereqRE.FindStringSubmatch(data)

	ci.Credits, _ = strconv.Atoi(prereqInt[1])
	ci.Course = ci.Course[0:3] + ci.Course[4:len(ci.Course)-1]
	ci.Prereqs = ParsePrereqs(ci.Course, prereqInt[2])

	return ci, nil
}

func GetClassInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	class := vars["class"]
	data, _ := GetClassPage(class)

	ci, _ := ParseClassPage(data)

	json.NewEncoder(w).Encode(ci)
}

func ServeFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving File: ", path.Clean(r.URL.RequestURI()))

	uri := "./templates/" + path.Clean(r.URL.RequestURI())
	http.ServeFile(w, r, uri)
}

func ServeCSS(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving File: ", path.Clean(r.URL.RequestURI()))

	uri := "./templates/css/" + path.Clean(r.URL.RequestURI())
	http.ServeFile(w, r, uri)
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api/{class}", GetClassInfo)
	router.PathPrefix("/").HandlerFunc(ServeFile)

	log.Fatal(http.ListenAndServe(":8080", router))
}
