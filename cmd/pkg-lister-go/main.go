package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ultrabluewolf/pkg-lister-go/cli"
	"github.com/ultrabluewolf/pkg-lister-go/stringarray"
)

type FilePackages []string

type Project struct {
	Path         string
	FilePackages map[string]FilePackages
}

func GetFilenames(path string) []string {
	ignoredDirs := []string{".git"}

	fileExtRE := regexp.MustCompile(`.*\.go`)

	filenames := []string{}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() && stringarray.Contains(ignoredDirs, info.Name()) {
			return filepath.SkipDir
		}

		if !info.IsDir() && fileExtRE.MatchString(info.Name()) {
			filenames = append(filenames, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return nil
	}

	return filenames
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
	f, err := os.Open(filename)
	PanicIf(err)
	defer f.Close()

	var (
		importSingleRE     = regexp.MustCompile(`^import "(.*)"$`)
		importGroupRE      = regexp.MustCompile("^import \\($")
		importGroupItemRE  = regexp.MustCompile(`"(.*)"$`)
		importGroupTermRE  = regexp.MustCompile("^\\)$")
		exhaustedImportsRE = regexp.MustCompile("^(func|type|const|var)")
	)

	scanner := bufio.NewScanner(f)
	pkgs := FilePackages{}
	isImportGroup := false
	for scanner.Scan() {
		line := scanner.Text()

		if !isImportGroup && importSingleRE.MatchString(line) {
			pkg := importSingleRE.FindStringSubmatch(line)[1]
			pkgs = append(pkgs, pkg)
			continue
		}

		if !isImportGroup && importGroupRE.MatchString(line) {
			isImportGroup = true
			continue
		}

		if isImportGroup && importGroupItemRE.MatchString(line) {
			pkg := importGroupItemRE.FindStringSubmatch(line)[1]
			pkgs = append(pkgs, pkg)
		}

		if isImportGroup && importGroupTermRE.MatchString(line) {
			isImportGroup = false
			continue
		}

		// imports are only valid near the top before declarations
		if exhaustedImportsRE.MatchString(line) {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

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

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func (project *Project) ToJSON() string {
	j, err := json.Marshal(project.ToPkgMapping())
	PanicIf(err)
	return string(j)
}

func main() {
	cmd := cli.ParseFlags()

	filenames := GetFilenames(cmd.ProjectPath)
	project := ExtractPackages(cmd.ProjectPath, filenames)
	output := project.ToJSON()

	fmt.Println(output)
}
