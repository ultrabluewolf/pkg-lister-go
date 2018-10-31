# pkg-lister-go

CLI that lists packages used by a specified project.

### Installation

`go get -u github.com/ultrabluewolf/pkg-lister-go/...`

Assuming your go bin is on your path to use the CLI see help for available commands:

`pkg-lister-go -help`

To display packages used by a project:

`pkg-lister-go -project <project-path>`

### Development

`go run cmd/pkg-lister-go/main.go -help|project <project-path>`

To display packages used by this project for example:

`go run cmd/pkg-lister-go/main.go -project .`
