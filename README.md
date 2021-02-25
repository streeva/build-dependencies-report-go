# Build Dependencies Report
Tooling that fetches detailed information for a list of dependencies (such as license information) from the package management ecosystem and builds an HTML report document.  The report will include a table of dependency information with the full detail of what was returned by the package management ecosystem included beneath it.

This tool is designed to work against the CSV file that is output by [streeva/read-dependencies-go](https://github.com/streeva/read-dependencies-go) or [streeva/read-dependencies-action](https://github.com/streeva/read-dependencies-action).

An exclude list can be provided if there are private project dependencies you don't want included in the report.

Group name is used for the document title, for handing in something like the repo name that gives the report a sensible context.

## Usage
A pre-built Docker image is available publicly on GitHub Container Registry, which you can run as so:
```
docker run -it -v `pwd`:/workspace -w /workspace ghcr.io/streeva/build-dependencies-report:v1.1.0 [parameters]
```
### Arguments
```bash
Usage of ./build-dependencies-report:
  -g string
    	Specify group name
  -i string
    	Specify input file name
  -o string
    	Specify report document name (default "dependency_report.html")
  -x string
    	Pattern for dependencies to exclude from the report
```

## Input CSV File format
```
<Source manifest file name>,<Package Management Ecosystem>,<Package Name>,<Package Version>
```
E.g.
```
streeva.csproj,NuGet,Microsoft.CodeAnalysis.CSharp,3.7.0
```

## Build
Clone the repo
```
git clone git@github.com:streeva/build-dependencies-report-go.git

cd build-dependencies-report-go
```
Build the application
```
go build
```
Run directly
```
./build-dependencies-report
```
Or build the Docker image
```
docker build . -t build-dependencies-report
```