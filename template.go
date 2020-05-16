package gomvc

import (
  "fmt"
  "log"
  "os"

  rice "github.com/GeertJohan/go.rice"
  "github.com/aymerick/raymond"
)

func defaultTemplateDir() string {
  path, err := os.Getwd()
  if err != nil {
    panic(err)
  }
  return fmt.Sprintf("%s/templates", path)
}

func createFileFromTemplates(template string, data interface{}, destPath string) error {
  box := rice.MustFindBox("templates")
  tmplString := box.MustString(template)
  tmpl, err := raymond.Parse(tmplString)
  if err != nil {
    return err
  }
  r := tmpl.MustExec(data)
  tmpl.RegisterHelper("whichAction", func(action string) string {
    log.Println("looking for HTTP action partial", action)
    if action == "" {
      log.Println("blank action name provided")
      return ""
    }
    return methodPartial(data, action, "gin")
  })
  if err := createFileFromString(destPath, r); err != nil {
    log.Println("could not create file for", destPath)
    return err
  }
  return nil
}

type TemplateHelper struct {
  Name     string
  Function func(string) string
}

func createFileWithHelpers(template string, data interface{}, destPath string, helpers []TemplateHelper) error {
  box := rice.MustFindBox("templates")
  tmplString := box.MustString(template)
  tmpl, err := raymond.Parse(tmplString)
  if err != nil {
    return err
  }
  for _, helper := range helpers {
    tmpl.RegisterHelper(helper.Name, helper.Function)
  }
  r := tmpl.MustExec(data)
  if err := createFileFromString(destPath, r); err != nil {
    log.Println("could not create file for", destPath)
    return err
  }
  return nil
}
