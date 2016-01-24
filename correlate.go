package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func contains(list []string, check string) bool {
	for _, item := range list {
		if strings.Contains(item, check) || strings.Contains(check, item) {
			return true
		}
	}

	return false
}

func removeSpaces(str string) string {
	return strings.Replace(str, " ", "", -1)
}

func checkRequirements(class ClassInfo, classesDone []string) (unmet []string) {
	unmet = make([]string, 0)

	for _, prereq := range class.Prereqs {
		if !contains(classesDone, prereq) {
			unmet = append(unmet, prereq)
		}
	}

	return unmet
}

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

type TemplateUnmetReq struct {
	Group       string
	CredsNeeded float64
	Classes     []ClassInfo
}

type TemplateClassList struct {
	Major   string
	College string
	Groups  []TemplateUnmetReq
}

func findClass(list []string, search string) int {
	for i, item := range list {
		if item[:7] == search[:7] {
			return i
		}
	}

	return -1
}

func CalcClassesJSON(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	unmet := r.FormValue("unmet")
	met := r.FormValue("met")

	fmt.Println("UNMET: ", unmet)
	fmt.Println("MET: ", met)

	aur := ParseUnmet(string(unmet))
	done := ParseAllClasses(string(met))

	byGroup := CalcAllClasses(aur, done)

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

	json.NewEncoder(w).Encode(tmplInfo)
}
