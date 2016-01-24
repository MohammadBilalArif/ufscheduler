package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func parseSection(data string) (hoursNeeded float64, group string, classes []string, remainder string) {
	lines := strings.Split(data, "\n")

	nextSection := -1

	if len(lines) == 0 {
		return
	}

	if len(lines[0]) == 0 {
		lines = lines[1:]
	}

	groupRE := regexp.MustCompile("Not Complete (?P<group>[A-Z \\-]*)=?")
	needsRE := regexp.MustCompile("NEEDS: (?P<needs>[0-9\\.]{2-4}) HOURS")
	classRE := regexp.MustCompile("(?P<class>[A-Z]{3}[0-9]{4}[A-Z]?)")

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if len(line) == 0 {
			nextSection = i
			break
		}

		gs := groupRE.FindStringSubmatch(line)
		ns := needsRE.FindStringSubmatch(line)
		cs := classRE.FindAllString(line, -1)

		if len(gs) == 2 {
			fmt.Println("GROUP: ", gs)
			group = gs[1]
		}

		if len(ns) == 2 {
			hoursNeeded, _ = strconv.ParseFloat(ns[1], 32)
		}

		if len(cs) > 0 {
			classes = append(classes, cs...)
		}
	}

	if nextSection == -1 {
		remainder = ""
	} else {
		remainder = strings.Join(lines[nextSection:], "\n")
	}

	return
}

func parseHeader(data string) (major, college, remainder string) {

	lines := strings.Split(data, "\n")

	beginning := -1
	state := 0

	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "Unmet Criteria") {
			state++

			for {
				if len(strings.TrimSpace(lines[i])) > 0 {
					i++
				} else {
					break
				}
			}

			continue
		}

		if len(strings.TrimSpace(lines[i])) == 0 {
			continue
		}

		switch state {
		case 1:
			if len(lines[i]) > 0 {
				major = strings.TrimSpace(lines[i])
				state++
			}
		case 2:
			if len(lines[i]) > 0 {
				college = strings.TrimSpace(lines[i])
				state++
			}
		case 3:
			beginning = i
			goto end
		}
	}

end:

	if beginning == -1 {
		log.Println("error: Invalid Format")

		remainder = ""
		return
	}

	remainder = strings.Join(lines[beginning:], "\n")

	return
}

type UnmetReq struct {
	Group       string
	Classes     []string
	HoursNeeded float64
}

type AllUnmetReqs struct {
	Major     string
	College   string
	UnmetReqs []UnmetReq
}

func ParseUnmet(wholeData string) (aur AllUnmetReqs) {
	var needs float64
	var group, major, college string
	var classes []string
	var data string

	major, college, data = parseHeader(string(wholeData))

	aur.Major = major
	aur.College = college

	fmt.Println("Major: ", aur.Major)
	fmt.Println("College: ", aur.College)

	for {
		needs, group, classes, data = parseSection(data)

		if len(group) == 0 || len(classes) == 0 {
			break
		}

		unmetReq := UnmetReq{
			Group:       group,
			HoursNeeded: needs,
			Classes:     classes,
		}

		fmt.Println("UnmetReq: ", unmetReq)

		aur.UnmetReqs = append(aur.UnmetReqs, unmetReq)

		if len(data) == 0 {
			break
		}
	}

	return
}

/*
func main() {
	file, _ := os.Open("unmet2.txt")
	defer file.Close()

	data, _ := ioutil.ReadAll(file)

	fmt.Println("Unmet: ", parseUnmet(string(data)))
}
*/
