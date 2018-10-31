package files

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ultrabluewolf/pkg-lister-go/stringarray"
)

func GetFilenames(projectPath string) []string {
	ignoredDirs := []string{".git"}
	fileExtRE := regexp.MustCompile(`.*\.go`)

	filenames := []string{}

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("failure accessing path %q: %v\n", path, err)
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
		fmt.Printf("failure walking path %q: %v\n", projectPath, err)
		return nil
	}

	return filenames
}
