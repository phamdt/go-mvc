package gomvc

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"

	"github.com/Fanatics/toast/collector"
)

// Seed is the cli command that creates new seeders based on structs in your
// application's "models" directory
var Seed = &cobra.Command{
	Use:   "seed",
	Short: "Generate seed files",
	Run: func(cmd *cobra.Command, args []string) {
		dest, _ := cmd.LocalFlags().GetString("dest")
		if dest == "" {
			path, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			dest = path
		}
		cmdDir := filepath.Join(dest, "cmd")
		createDirIfNotExists(cmdDir)
		if err := CreateSeedCommand(dest); err != nil {
			panic(err)
		}
	},
}

// CreateSeedCommand creates a single cmd that inserts records for all models
func CreateSeedCommand(dest string) error {
	modelDir := filepath.Join(dest, "models")
	data, err := collectFileData(modelDir)
	if err != nil {
		log.Println("filepath walk error", err)
		return err
	}
	var modelNames []string
	iterateModels(data, func(s collector.Struct) error {
		modelNames = append(modelNames, s.Name)
		return nil
	})
	if len(modelNames) == 0 {
		return errors.New("No models found")
	}
	content := createContentsFromTemplate("sqlboiler/seeder.tmpl", struct {
		Models []string
	}{Models: modelNames})
	commandDir := filepath.Join(dest, "cmd/seed")
	createDirIfNotExists(commandDir)
	fileName := filepath.Join(commandDir, "main.go")
	return createFileFromString(fileName, content)
}

// CreateSeederWithName creates a file with a single function. This function
// connects to the DB and inserts records for each model using each the model
// factories.
func CreateSeederWithName(structName string, destDir string) error {
	// default to current dir if none given
	if destDir == "" {
		destDir = "./cmd/seed"
	}
	createDirIfNotExists(destDir)
	fileName := fmt.Sprintf("%s.go", strcase.ToSnake(structName))
	dest := filepath.Join(destDir, fileName)
	name := strcase.ToCamel(structName)
	data := map[string]string{"Name": name}
	contents := createContentsFromTemplate("sqlboiler/seed.tmpl", data)
	return createFileFromString(dest, contents)
}

// CreateSeedersFromModels iterates through the given directory
// finds models the creates seed files in the given destination
func CreateSeedersFromModels(dir string, dest string) error {
	data, err := collectFileData(dir)
	if err != nil {
		log.Println("filepath walk error", err)
		return err
	}
	log.Println(len(data.Packages))
	if len(data.Packages) == 0 {
		log.Println("No model data found. Exiting.")
		return nil
	}
	iterateModels(data, func(s collector.Struct) error {
		if err := CreateSeederWithName(s.Name, dest); err != nil {
			return err
		}
		return nil
	})

	return nil
}

func iterateModels(data *collector.Data, lambda func(s collector.Struct) error) {
	for _, p := range data.Packages {
		for _, f := range p.Files {
			for _, s := range f.Structs {
				var hasAnyTags bool // hack
				for _, field := range s.Fields {
					if len(field.Tag) > 0 {
						hasAnyTags = true
						tag := strings.Trim(field.Tag, "`")
						parts := strings.Split(tag, " ")
						for _, p := range parts {
							tag := strings.Split(p, ":")
							name := tag[0]
							// some db tools use a 'db' struct tag to help with marshalling
							// so we use that convention here as a way to determine what the sql column is
							if name == "db" {
								values := strings.Trim(tag[1], "\"")
								valueParts := strings.Split(values, ",")
								log.Println(valueParts)
							}
						}
					}
				}
				if !hasAnyTags {
					continue
				}
				lambda(s)
			}
		}
	}
}

func collectFileData(dir string) (*collector.Data, error) {
	fset := token.NewFileSet()
	data := &collector.Data{}
	err := filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Fatal("recursive walk error:", err)
		}
		// skip over files, only continue into directories for parser to enter
		if !fi.IsDir() {
			return nil
		}

		pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Fatalf("parse dir error: %v\n", err)
		}
		for _, pkg := range pkgs {
			p := collector.Package{
				Name: pkg.Name,
			}
			for _, file := range pkg.Files {
				c := &collector.FileCollector{}
				ast.Walk(c, file)
				f := collector.File{
					Name:             fset.Position(file.Pos()).Filename,
					Package:          pkg.Name,
					Imports:          c.Imports,
					BuildTags:        c.BuildTags,
					Comments:         c.Comments,
					MagicComments:    c.MagicComments,
					GenerateComments: c.GenerateComments,
					Consts:           c.Consts,
					Vars:             c.Vars,
					Structs:          c.Structs,
					TypeDefs:         c.TypeDefs,
					Interfaces:       c.Interfaces,
					Funcs:            c.Funcs,
				}
				p.Files = append(p.Files, f)
			}
			data.Packages = append(data.Packages, p)
		}

		return nil
	})
	return data, err
}
