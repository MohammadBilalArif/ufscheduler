package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	UsersHelped    = 0
	TotalTimeTaken = 0
)

// contains checks if a string is contained within another string
func contains(list []string, check string) bool {
	for _, item := range list {
		if strings.Contains(item, check) || strings.Contains(check, item) {
			return true
		}
	}

	return false
}

// removeSpaces removes all spaces from a string (including middle ones)
func removeSpaces(str string) string {
	return strings.Replace(str, " ", "", -1)
}

// checkRequirements checks the list of classesDone for all the classes in
// class.Prereqs, returning a list of all classes that have not been completed
func checkRequirements(class ClassInfo, classesDone []string) (unmet []string) {
	unmet = make([]string, 0)

	for _, prereq := range class.Prereqs {
		if !contains(classesDone, prereq) {
			unmet = append(unmet, prereq)
		}
	}

	return unmet
}

// calcNeededClasses calculates the classes that can be taken based on prerequisites, as well as generating
// a list of classes that have NOT been completed (prereqs)
func calcNeededClasses(neededClasses []string, doneClasses []string) (ready []string, prereqs []string) {
	needed := make([]string, 0)
	good2go := make([]string, 0)

	for _, class := range neededClasses {
		class = removeSpaces(class)

		info, err := GetClassInfo(class)
		if err != nil {
			fmt.Printf("ERROR ON CLASS: '%s'\n", class)
			continue
		}

		reqs := checkRequirements(info, doneClasses)

		if len(reqs) > 0 {
			good2go = append(good2go, class)
			needed = append(needed, reqs...)
		} else {
			good2go = append(good2go, class)
		}
	}

	ready = good2go
	prereqs = needed

	return
}

// removeDuplicates removes all duplicates from a list of classes
func removeDuplicates(list []string) (result []string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range list {
		if !found[x] {
			found[x] = true
			list[j] = list[i]
			j++
		}
	}
	list = list[:j]

	return list
}

// removeReverseDuplicates attempts to remove all duplicates starting from the
// end and moving towards the front
func removeReverseDuplicates(list []string) []string {
	reverse := make([]string, len(list))

	for i, item := range list {
		reverse[len(list)-1-i] = item
	}

	list = removeDuplicates(list)

	reverse = make([]string, len(list))

	for i, item := range list {
		reverse[len(list)-1-i] = item
	}

	return reverse
}

// CalcAllClasses is the main function that wraps all the other helper functions in
// determining which classes should be taken in which order.
func CalcAllClasses(aur AllUnmetReqs, doneClasses []string) (byGroup [][]string) {
	goodList := make([]string, 0)
	prereqList := make([]string, 0)
	doneClasses = removeReverseDuplicates(doneClasses)

	byGroup = make([][]string, 0)

	for _, unmet := range aur.UnmetReqs {

		prereqList = make([]string, 0)
		goodList = make([]string, 0)

		for _, class := range unmet.Classes {
			prereqList = append(prereqList, class)
		}

		for i := 0; i < 15; i++ { // set maximum 15 semesters
			g, p := calcNeededClasses(prereqList, doneClasses)

			goodList = append(g, goodList...)
			prereqList = p

			if len(prereqList) == 0 {
				break
			}
		}

		byGroup = append(byGroup, removeReverseDuplicates(goodList))
	}

	return byGroup
}

// TemplateUnmetReq is an unmet requirement in JSON friendly format
type TemplateUnmetReq struct {
	Group       string
	CredsNeeded float64
	Classes     []ClassInfo
}

// TemplateClassList is a class list in JSON friendly format
type TemplateClassList struct {
	Major   string
	College string
	Groups  []TemplateUnmetReq
}

// findClass searches a list of classes for a particular one,
// allowing for leeway when classes may or may not have a letter
// after them (the UF website is bad about this!)
func findClass(list []string, search string) int {
	for i, item := range list {
		if item[:7] == search[:7] {
			return i
		}
	}

	return -1
}

// CalcClassesJSON wraps CalcAllClasses and returns the results
// in JSON format, useful for our website (and any other REST endpoint users)
func CalcClassesJSON(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	defer func() {
		endTime := time.Now()

		TotalTimeTaken += int(endTime.Sub(startTime).Nanoseconds() / 1e6)
	}()

	r.ParseForm()

	unmet := r.FormValue("unmet")
	met := r.FormValue("met")

	log.Printf("Met: %s", met)
	log.Printf("Unmet: %s", unmet)

	UsersHelped++

	log.Println("Parsing Unmet Requirements")

	aur := ParseUnmet(string(unmet))

	log.Println("Parsing All Classes")

	done := ParseAllClasses(string(met))

	log.Println("Calculating Intersection of Classes")

	byGroup := CalcAllClasses(aur, done)

	log.Println("Double-checking Class Ordering")

	// this loop double checks that all classes prerequisites are met before a class
	// can be taken.  Our heuristic is imperfect and this was a quick solution.
	for j := 0; j < 3; j++ {
		for k, group := range byGroup {
			for i, class := range group {
				info, err := GetClassInfo(class)
				if info.Course == "" {
					continue
				}
				if err != nil {
					continue
				}

				for _, prereq := range info.Prereqs {
					index := findClass(group, prereq)

					if index > i {
						group[i], group[index] = group[index], group[i]
					}
				}

				byGroup[k] = group
			}
		}
	}

	log.Println("Preparing Data for JSON Format")

	// generate the data in a format that is JSON friendly
	infos := make([]TemplateUnmetReq, 0)

	for i, group := range byGroup {
		grp := make([]ClassInfo, 0)
		unmetReq := aur.UnmetReqs[i]

		for _, class := range group {
			info, err := GetClassInfo(class)
			if info.Course == "" {
				continue
			}
			if err != nil {
				fmt.Println("error finding class", class)
				continue
			}

			grp = append(grp, info)
		}

		infos = append(infos, TemplateUnmetReq{
			Group:       unmetReq.Group,
			CredsNeeded: unmetReq.HoursNeeded,
			Classes:     grp,
		})
	}

	tmplInfo := TemplateClassList{
		Major:   aur.Major,
		College: aur.College,
		Groups:  infos,
	}

	w.WriteHeader(200)

	log.Println("Sending Back JSON Data")

	// send out all our data using the JSON encoder
	json.NewEncoder(w).Encode(tmplInfo)
}
