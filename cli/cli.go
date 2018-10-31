package cli

import (
	"flag"
	"fmt"
	"os"
)

type CmdOpts struct {
	ProjectPath string
}

func GetUsageText() string {
	return fmt.Sprintf(
		"[pkg-lister-go] A CLI that lists packages used by a specified Go project.\n\n%s%s",
		"USAGE:\n\n",
		"   pkg-lister-go -project <project-path>\n",
	)
}

func ParseFlags() *CmdOpts {
	var project string
	flag.StringVar(&project, "project", "", "Go project")

	var help bool
	flag.BoolVar(&help, "help", false, "Help")
	flag.Parse()

	exitCode := 0
	if !help && len(project) == 0 {
		exitCode = 1
	}

	if help || len(project) == 0 {
		usageTxt := GetUsageText()
		fmt.Println(usageTxt)
		os.Exit(exitCode)
	}

	return &CmdOpts{
		ProjectPath: project,
	}
}
