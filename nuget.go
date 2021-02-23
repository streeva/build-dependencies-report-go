package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kyoh86/go-spdx/spdx"
)

type Package struct {
	Metadata Metadata `xml:"metadata"`
}

type Metadata struct {
	Id												string				`xml:"id"`
	Version										string				`xml:"version"`
	Title											string				`xml:"title"`
	Authors										string				`xml:"authors"`
	Owners										string				`xml:"owners"`
	DevelopmentDependency			bool					`xml:"developmentDependency"`
	License										SPDXLicense		`xml:"license"`
	LicenseUrl								string				`xml:"licenseUrl"`
	ProjectUrl								string				`xml:"projectUrl"`
	Description								string				`xml:"description"`
	Repository								Repository		`xml:"repository"`
	Dependencies							Dependencies	`xml:"dependencies"`
}

type SPDXLicense struct {
	Identifier	string `xml:",chardata"`
	Type				string `xml:"type,attr"`
}

type Repository struct {
	Type		string	`xml:"type,attr"`
	Url			string	`xml:"url,attr"`
	Branch	string	`xml:"branch,attr"`
	Commit	string	`xml:"commit,attr"`
}

type Dependencies struct {
	Dependencies	[]PackageDependency	`xml:"dependency"`
	Groups				[]Group							`xml:"group"`
}

type Group struct {
	TargetFramework	string							`xml:"targetFramework,attr"`
	Dependencies		[]PackageDependency	`xml:"dependency"`
}

type PackageDependency struct {
	Id			string	`xml:"id,attr"`
	Version	string	`xml:"version,attr"`
}

func ReadDependencyExtInfo(dependencyMap map[Dependency]DependencyExtInfo) error {
	for dependency := range dependencyMap {
		dependencyExtInfo, err := ReadExtInfoFromNuGet(dependency)
		if err != nil {
			return err
		}

		dependencyMap[dependency] = dependencyExtInfo
	}

	return nil
}

func ReadExtInfoFromNuGet(dependency Dependency) (DependencyExtInfo, error) {
	client := &http.Client{ Timeout: time.Second * 10 }
	response, err := client.Get(fmt.Sprintf("https://api.nuget.org/v3-flatcontainer/%[1]s/%[2]s/%[1]s.nuspec", dependency.Name, dependency.Version))
	if err != nil{
		return DependencyExtInfo{}, err
	}

	bodyText, err := ioutil.ReadAll(response.Body)
	if err != nil{
		return DependencyExtInfo{}, err
	}

	var pkg Package
	xml.Unmarshal(bodyText, &pkg)
	license, licenseLink, _ := ResolveLicense(pkg.Metadata)
	owners := pkg.Metadata.Owners
	if len(owners) <= 0 {
		owners = pkg.Metadata.Authors
	}

	var dependencyExtInfo = DependencyExtInfo{
		Owners: owners,
		ProjectUrl: pkg.Metadata.ProjectUrl,
		License: License{ license, licenseLink },
		Description: pkg.Metadata.Description,
		Dependencies: make([]string, 0),
		DevelopmentDependency: pkg.Metadata.DevelopmentDependency,
		Raw: string(bodyText),
	}

	return dependencyExtInfo, nil
}

func ResolveLicense(metadata Metadata) (string, string, error) {
	var licenseName string
	var licenseUrl	string
	if metadata.License.Type == "expression" {
		tree, err := spdx.Parse(metadata.License.Identifier)
		if err != nil {
			return "", "", err
		}

		// This will probably break if we get any packages with a compound identifier
		l, err := spdx.Get(tree.String())
		if err != nil {
			return "", "", err
		}
		
		licenseName = l.Name
		licenseUrl = l.URL
	} else if len(metadata.LicenseUrl) > 0 {
		licenseName = "License Link"
		licenseUrl = metadata.LicenseUrl
	} else {
		licenseName = "Project Link"
		licenseUrl = metadata.ProjectUrl
	}

	return licenseName, licenseUrl, nil
}