package main

import (
	"fmt"
	"regexp"

	"github.com/ultrabluewolf/pkg-lister-go/cli"
	"github.com/ultrabluewolf/pkg-lister-go/files"
	"github.com/ultrabluewolf/pkg-lister-go/project"
)

func main() {
	// filename filtering settings
	ignoredDirs := []string{".git"}
	fileExtRE := regexp.MustCompile(`.*\.go`)

	cmd := cli.ParseFlags()

	filenames := files.GetFilenames(cmd.ProjectPath, fileExtRE, ignoredDirs)
	proj := project.ExtractPackages(cmd.ProjectPath, filenames)
	output := proj.ToJSON()

	fmt.Println(output)
}
