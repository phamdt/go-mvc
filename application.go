package gomvc

import (
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/GeertJohan/go.rice/embedded"
	"github.com/spf13/cobra"
)

// Application is the cli command that creates new application
var Application = &cobra.Command{
	Use:   "application",
	Short: "Generate application files",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a name for your new application")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		log.Printf("preparing to create a new application: %s\n", appName)

		destinationDir := getAppDir(cmd, appName)
		// setup basic directories
		createDirIfNotExists(destinationDir)
		createDirIfNotExists(path.Join(destinationDir, "controllers"))
		createDirIfNotExists(path.Join(destinationDir, "models"))
		createDirIfNotExists(path.Join(destinationDir, "migrations"))
		createDirIfNotExists(path.Join(destinationDir, ".circleci"))

		// copy over static go files (not templated)
		for filename := range embedded.EmbeddedBoxes["static"].Files {
			log.Println("copying static file", destinationDir, filename)
			copyStatic(destinationDir, filename)
		}
		log.Println("finished copying static files")
		// render files from generic gomvc templates
		for _, file := range []File{
			{Template: "sqlboiler/query.go.tpl", Name: "controllers/query.go"},
			{Template: "gin/main.tpl", Name: "main.go"},
			{Template: "build/docker-compose.yml.tpl", Name: "docker-compose.yml"},
			{Template: "build/env.tpl", Name: ".env"},
			{Template: "build/wait-for-server-start.sh.tpl", Name: ".circleci/wait-for-server-start.sh"},
		} {
			data := map[string]string{
				"Name":      appName,
				"TitleName": strings.Title(appName),
			}
			destPath := filepath.Join(destinationDir, file.Name)
			createFileFromTemplates(file.Template, data, destPath)
		}
		// render files from special gomvc templates with specific template data
		ctrlDir := filepath.Join(destinationDir, "controllers")
		CreateRouter(RouteData{}, "gin/router.tpl", ctrlDir)
	},
	PostRun: func(cmd *cobra.Command, args []string) { // this doesn't work for some reason
		appName := args[0]
		appDir := getAppDir(cmd, appName)

		// gofmt
		log.Println("running gofmt on", appDir)
		runGoFmt(appDir)

		// goimports
		log.Println("running goimports on", appDir)
		runGoImports(appDir)

		// go module
		log.Println("creating go module", appName)
		createModule(appDir, appName)
	},
}

func getAppDir(cmd *cobra.Command, appName string) string {
	dest, err := cmd.LocalFlags().GetString("dest")
	if err != nil {
		panic(err)
	}
	if dest == "" {
		cwd, _ := os.Getwd()
		return path.Join(cwd, appName)
	}
	return dest
}

func runGoFmt(appDir string) {
	command := exec.Command("gofmt", "-w", appDir)
	runCommand(command)
	log.Printf("Just ran gofmt subprocess %d, exiting\n", command.Process.Pid)
}

func runGoImports(appDir string) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	command := exec.Command(fmt.Sprintf("%s/bin/goimports", gopath), "-w", appDir)
	runCommand(command)
	log.Printf("Just ran goimports subprocess %d, exiting\n", command.Process.Pid)
}

// currently can only be used in app dir
func createModule(appDir, appName string) {
	command := exec.Command("go", "mod", "init", appName)
	command.Dir = appDir
	runCommand(command)
}

func runCommand(command *exec.Cmd) {
	stderr, err := command.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := command.Start(); err != nil {
		log.Fatal(err)
	}
	slurp, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", slurp)
	if err := command.Wait(); err != nil {
		log.Fatal(err)
	}
}

// File is a GoMVC specific type to store rendering meta data with the filenames
type File struct {
	Template string
	Name     string
}

func copyStatic(destinationBasePath string, name string) {
	box := rice.MustFindBox("static")
	dest := filepath.Join(destinationBasePath, name)
	if err := createFileFromString(dest, box.MustString(name)); err != nil {
		panic(err)
	}
}
