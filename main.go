package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"regexp"
)

type Dependency struct {
	Name    string
	Version string
}

type License struct {
	Type string
	Link string
}

type DependencyExtInfo struct {
	ProjectUrl            string
	Repository            string
	License               License
	Description           string
	Owners                string
	Dependencies          []string
	DevelopmentDependency bool
	Raw                   string
}

const EXIT_SUCCESS = 0
const EXIT_FAILURE = 1

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(EXIT_FAILURE)
	}
}

func main() {
	var filename string
	var groupName string
	var reportName string
	var exclude string
	flag.StringVar(&filename, "i", "", "Specify input file name")
	flag.StringVar(&groupName, "g", "", "Specify group name")
	flag.StringVar(&reportName, "o", "dependency_report.html", "Specify report document name")
	flag.StringVar(&exclude, "x", "", "Pattern for dependencies to exclude from the report")
	flag.Parse()

	if len(filename) <= 0 || len(groupName) <= 0 {
		fmt.Println("Please specify filename containing the dependency data and a group name to report the projects under")
		os.Exit(EXIT_FAILURE)
	}

	// De-duplicated list of dependencies to look up more information on
	dependencyTree := make(map[string]map[Dependency]DependencyExtInfo)
	// Dependencies by project name for organising the report
	usageTree := make(map[string][]Dependency)
	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	check(err)

	for _, record := range records {
		if len(record) < 4 {
			fmt.Println("Unexpected line in input file")
			os.Exit(EXIT_FAILURE)
		}

		dependency := Dependency{Name: record[2], Version: record[3]}
		if matched, _ := regexp.MatchString(exclude, dependency.Name); len(exclude) > 0 && matched {
			continue
		}

		if _, ok := dependencyTree[record[1]]; !ok {
			dependencyTree[record[1]] = make(map[Dependency]DependencyExtInfo)
		}

		dependencyTree[record[1]][dependency] = DependencyExtInfo{}
		usageTree[record[0]] = append(usageTree[record[0]], dependency)
	}

	for ecosystem, dependencyMap := range dependencyTree {
		fmt.Printf("Processing ecosystem %s\n", ecosystem)
		err := ReadDependencyExtInfo(dependencyMap)
		check(err)
	}

	_ = GenerateReport(reportName, "Dependency Report for "+groupName, usageTree, dependencyTree)
}
