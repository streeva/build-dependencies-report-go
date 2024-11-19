package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Table struct {
	Rows []TableRow
}

type TableRow struct {
	Project    string
	ProjectUrl string
	Dependency Dependency
	Ecosystem  string
	License    License
}

type alphabetically []TableRow

func (s alphabetically) Len() int {
	return len(s)
}
func (s alphabetically) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s alphabetically) Less(i, j int) bool {
	if s[i].Project == s[j].Project {
		if s[i].Dependency.Name == s[j].Dependency.Name {
			return s[i].Dependency.Version < s[j].Dependency.Version
		}

		return s[i].Dependency.Name < s[j].Dependency.Name
	}

	return s[i].Project < s[j].Project
}

type PackageDetail struct {
	EcosystemName string
	Dependencies  map[Dependency]DependencyExtInfo
}

func GenerateReport(filename string, documentTitle string, usageInfo map[string][]Dependency, dependencyInfo map[string]map[Dependency]DependencyExtInfo) error {

	var table Table
	// Compile dependency information by project
	for project, dependencies := range usageInfo {
		for _, dependency := range dependencies {
			const ecosystem = "NuGet" // TODO: slightly cheeky assumption, to revisit when I add other ecosystems
			dependencyExtInfo := dependencyInfo[ecosystem][dependency]
			table.AddRow(&TableRow{project, dependencyExtInfo.ProjectUrl, dependency, ecosystem, dependencyExtInfo.License})
		}
	}

	// Report detailed information for each dependency
	var packageDetails []PackageDetail
	for ecosystem, dependencies := range dependencyInfo {
		packageDetails = append(packageDetails, PackageDetail{ecosystem, dependencies})
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	writeHtmlHeader(w, documentTitle)

	for _, line := range table.GetHtmlLines() {
		w.WriteString(line + "\n")
	}

	for _, packageDetail := range packageDetails {
		for _, line := range packageDetail.GetHtmlLines() {
			w.WriteString(line + "\n")
		}
	}

	writeHtmlFooter(w)
	return nil
}

func (table *Table) AddRow(row *TableRow) {
	table.Rows = append(table.Rows, *row)
}

func (table *Table) GetHtmlLines() []string {
	var tableLines []string
	tableLines = append(tableLines, "<table>")
	tableLines = append(tableLines, "<tr><th>Project</th><th>Ecosystem</th><th>Library</th><th>Version</th><th>License</th><th>Raw</th></tr>")
	sort.Sort(alphabetically(table.Rows))
	for _, row := range table.Rows {
		tableLines = append(tableLines, row.BuildHtml())
	}
	tableLines = append(tableLines, "</table>")
	return tableLines
}

func (row *TableRow) BuildHtml() string {
	tableRow := "<tr>"
	tableRow += "<td>" + fileNameWithoutExtension(row.Project) + "</td>"
	tableRow += "<td>" + row.Ecosystem + "</td>"
	tableRow += fmt.Sprintf("<td><a href=\"%s\">%s</a></td>", row.ProjectUrl, row.Dependency.Name)
	tableRow += "<td>" + row.Dependency.Version + "</td>"
	tableRow += fmt.Sprintf("<td><a href=\"%s\">%s</a></td>", row.License.Link, row.License.Type)
	tableRow += fmt.Sprintf("<td><a href=\"#%s\">Detail</a></td>", row.Dependency.GetReference())
	tableRow += "</tr>"
	return tableRow
}

func (pkg *PackageDetail) GetHtmlLines() []string {
	var lines []string
	lines = append(lines, "<h2>Package Details from "+pkg.EcosystemName+"<h2>")
	for dependency, dependencyExtInfo := range pkg.Dependencies {
		detailLine := fmt.Sprintf("<h3><a id=\"%[1]s\"></a>%[1]s</h3><pre>%[2]s</pre>", dependency.GetReference(), escape(dependencyExtInfo.Raw))
		lines = append(lines, detailLine)
	}
	return lines
}

func (dep *Dependency) GetReference() string {
	return dep.Name + "@" + dep.Version
}

func writeHtmlHeader(writer *bufio.Writer, documentTitle string) {
	writer.WriteString("<html><head><title>" + documentTitle + "</title></head>\n")
	writer.WriteString("<body>\n")
	writer.WriteString("<h2>" + documentTitle + "</h2>\n")
}

func writeHtmlFooter(writer *bufio.Writer) {
	writer.WriteString("<footer><p align=\"center\">Report generated " + time.Now().Format("2006-01-02 15:04:05") + "</p></footer>")
	writer.WriteString("</body></html>\n")
}

func escape(input string) string {
	input = strings.ReplaceAll(input, "&", "&amp;")
	input = strings.ReplaceAll(input, "<", "&lt;")
	return input
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
