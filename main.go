package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Dependency struct {
	Name		string
	Version	string
}

type License struct {
	Type	string
	Link	string
}

type DependencyExtInfo struct {
	ProjectUrl								string
	Repository								string
	License										License
	Description								string
	Owners										string
	Dependencies							[]string
	DevelopmentDependency			bool
	Raw												string
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
	var groupname string
	var reportname string
	flag.StringVar(&filename, "i", "", "Specify input file name")
	flag.StringVar(&groupname, "g", "", "Specify group name")
	flag.StringVar(&reportname, "o", "dependency_report.md", "Specify report document name")
	flag.Parse()

	if len(filename) <= 0 || len(groupname) <= 0 {
		fmt.Println("Please specify filename containing the dependency data and a group name to report the projects under")
		os.Exit(EXIT_FAILURE)
	}

	// De-duplicated list of dependencies to look up more information on
	dependencyTree := make(map[string]map[Dependency]DependencyExtInfo)
	// Record which projects use which dependencies for later grouping in the report
	usageTree := make(map[string][]Dependency)
	file, err := os.Open(filename)
	check(err)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var segments = strings.Split(scanner.Text(), ",")
		if len(segments) < 4 {
			fmt.Println("Unexpected line in input file")
			os.Exit(EXIT_FAILURE)
		}

		dependency := Dependency{ Name: segments[2], Version: segments[3] }
		// Skip internal package use - TODO: more generic exclude list rather than hard-coding streeva
		if strings.Contains(dependency.Name, "streeva") {
			continue
		}

		if _, ok := dependencyTree[segments[1]]; !ok {
			dependencyTree[segments[1]] = make(map[Dependency]DependencyExtInfo)
		}

		dependencyTree[segments[1]][dependency] = DependencyExtInfo{}
		usageTree[segments[0]] = append(usageTree[segments[0]], dependency)
	}

	for ecosystem, dependencyMap := range dependencyTree {
		fmt.Printf("Processing ecosystem %s\n", ecosystem)
		err := ReadDependencyExtInfo(dependencyMap)
		check(err)
	}

	_ = GenerateReport(reportname, groupname, usageTree, dependencyTree)
}