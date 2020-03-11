package gomvc

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
)

func createFileFromString(filepath string, contents string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return errors.Wrap(err, "os.Create error")
	}
	w := bufio.NewWriter(f)
	_, err = w.WriteString(contents)
	w.Flush()

	if err != nil {
		return errors.Wrap(err, "write string error")
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

// TODO, return error?
func getFilesFromDir(dir string) []string {
	fileNames := []string{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
		return fileNames
	}
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames
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
		os.Mkdir(dir, os.ModePerm)
	}
}

func dirExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}
