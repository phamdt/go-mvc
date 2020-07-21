package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aymerick/raymond"
	gomvc "github.com/go-generation/go-mvc"
	"github.com/spf13/cobra"
)

func main() {
	root := Root()
	root.AddCommand(Command())

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var root = &cobra.Command{
	Use:   "dev",
	Short: "dev is the gomvc development helper",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
	},
}

func Root() *cobra.Command {
	return root
}

var command = &cobra.Command{
	Use:   "command",
	Short: "Generate cli commands",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a name for your new command")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := strings.ToLower(args[0])
		log.Println("preparing to create a new command:", name)

		data := map[string]string{
			"Name":      name,
			"TitleName": strings.Title(name),
		}
		t, err := raymond.ParseFile("./cmd/dev/command.tpl")
		if err != nil {
			panic(err)
		}
		r, err := t.Exec(data)
		if err != nil {
			panic(err)
		}
		destPath := fmt.Sprintf("%s/%s.go", ".", name)
		if err := gomvc.CreateFileFromString(destPath, r); err != nil {
			panic(err)
		}
	},
}

// Command is the cli command that creates new cli commands
func Command() *cobra.Command {
	return command
}
