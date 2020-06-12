package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	gomvc "github.com/go-generation/go-mvc"
	"github.com/hoisie/mustache"
	"github.com/spf13/cobra"
)

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

		templateDir := "./templates"
		commandDir := fmt.Sprintf("%s/command.tpl", templateDir)
		data := map[string]string{
			"Name":      name,
			"TitleName": strings.Title(name),
		}
		r := mustache.RenderFile(commandDir, data)
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
