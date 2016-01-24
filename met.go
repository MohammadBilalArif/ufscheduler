package main

import (
	"regexp"
	"strings"
)

func ParseAllClasses(data string) (classes []string) {
	classRE := regexp.MustCompile("[0-9][ \t]+(?P<class>[A-Z]{3}[0-9]{4}[A-Z]?)")
	lines := strings.Split(data, "\n")

	classes = make([]string, 0)

	for i := 0; i < len(lines); i++ {
		cs := classRE.FindAllStringSubmatch(lines[i], -1)

		for j := 0; j < len(cs); j++ {
			classes = append(classes, cs[j][1])
		}
	}

	return
}

/*
func main() {
	file, _ := os.Open("met2.txt")
	defer file.Close()

	data, _ := ioutil.ReadAll(file)

	fmt.Println("Classes: ", parseAllClasses(string(data)))
}
*/
