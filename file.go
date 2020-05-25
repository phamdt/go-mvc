package gomvc

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	rice "github.com/GeertJohan/go.rice"
	"github.com/aymerick/raymond"
	"github.com/pkg/errors"
)

func createContentsFromTemplate(tmplPath string, data interface{}) string {
	box := rice.MustFindBox("templates")
	tmplString := box.MustString(tmplPath)
	tmpl, err := raymond.Parse(tmplString)
	if err != nil {
		panic(err)
	}
	result := tmpl.MustExec(data)
	return result
}

// createFileFromString takes a filepath as the destination of the file
// to be created as well as the contents to be written to this file.
func createFileFromString(filepath string, contents string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return errors.Wrap(err, "createFileFromString: os.Create error")
	}
	w := bufio.NewWriter(f)
	_, err = w.WriteString(contents)
	w.Flush()

	if err != nil {
		return errors.Wrap(err, "createFileFromString: write string error")
	}
	return nil
}

func createStringFromFile(filePath string) string {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
// https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func createDirIfNotExists(dir string) {
	if !dirExists(dir) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			panic(err)
		}
		log.Printf("created %s\n", dir)
	}
}

func dirExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println(err.Error())
		return false
	}
	return true
}

func addGoExt(s string) string {
	return fmt.Sprintf("%s.go", s)
}
