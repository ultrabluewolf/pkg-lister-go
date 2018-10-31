package project

import (
	"encoding/json"
	"regexp"

	"github.com/ultrabluewolf/pkg-lister-go/exit"
	"github.com/ultrabluewolf/pkg-lister-go/files"
)

type FilePackages []string

type Project struct {
	Path         string
	FilePackages map[string]FilePackages
}

func ExtractPackages(path string, filenames []string) *Project {
	project := Project{
		Path:         path,
		FilePackages: map[string]FilePackages{},
	}

	for _, filename := range filenames {
		pkgs := ExtractPackagesForFilename(filename)
		project.FilePackages[filename] = pkgs
	}
	return &project
}

func ExtractPackagesForFilename(filename string) FilePackages {
	var (
		importSingleRE     = regexp.MustCompile(`^import "(.*)"$`)
		importGroupRE      = regexp.MustCompile("^import \\($")
		importGroupItemRE  = regexp.MustCompile(`"(.*)"$`)
		importGroupTermRE  = regexp.MustCompile("^\\)$")
		exhaustedImportsRE = regexp.MustCompile("^(func|type|const|var)")
	)

	pkgs := FilePackages{}
	isImportGroup := false

	processLine := func(line string) {
		if !isImportGroup && importSingleRE.MatchString(line) {
			pkg := importSingleRE.FindStringSubmatch(line)[1]
			pkgs = append(pkgs, pkg)
			return
		}

		if !isImportGroup && importGroupRE.MatchString(line) {
			isImportGroup = true
			return
		}

		if isImportGroup && importGroupItemRE.MatchString(line) {
			pkg := importGroupItemRE.FindStringSubmatch(line)[1]
			pkgs = append(pkgs, pkg)
		}

		if isImportGroup && importGroupTermRE.MatchString(line) {
			isImportGroup = false
			return
		}

	}

	// imports are only valid near the top before declarations
	shouldTerminate := func(line string) bool {
		return exhaustedImportsRE.MatchString(line)
	}

	files.ReadFileForEachLine(filename, processLine, shouldTerminate)

	return pkgs
}

func (project *Project) ToPkgMapping() map[string][]string {
	m := map[string][]string{}
	for filename, pkgs := range project.FilePackages {
		for _, pkg := range pkgs {
			if _, ok := m[pkg]; !ok {
				m[pkg] = []string{}
			}

			m[pkg] = append(m[pkg], filename)
		}
	}
	return m

}

func (project *Project) ToJSON() string {
	j, err := json.Marshal(project.ToPkgMapping())
	exit.PanicIf(err)
	return string(j)
}
