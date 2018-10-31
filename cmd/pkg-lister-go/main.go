package main

import (
	"fmt"

	"github.com/ultrabluewolf/pkg-lister-go/cli"
	"github.com/ultrabluewolf/pkg-lister-go/files"
	"github.com/ultrabluewolf/pkg-lister-go/project"
)

func main() {
	cmd := cli.ParseFlags()

	filenames := files.GetFilenames(cmd.ProjectPath)
	proj := project.ExtractPackages(cmd.ProjectPath, filenames)
	output := proj.ToJSON()

	fmt.Println(output)
}
