package gomvc

import (
	"bytes"
	"fmt"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

// CreateRouter creates a router.go file to be used in the controllers directory
func CreateRouter(data RouteData, relativeTemplatePath, destDir string) {
	sort.Slice(data.Controllers, func(i, j int) bool {
		return strings.Compare(
			data.Controllers[i].Name,
			data.Controllers[j].Name) < 1
	})
	outputPath := filepath.Join(destDir, "router.go")
	if err := createFileFromTemplates(relativeTemplatePath, data, outputPath); err != nil {
		log.Fatalf("error generating router file for %s %s\n", outputPath, err.Error())
	}
}

func AddActionViaAST(data ControllerData, routerFilePath string, destDir string) {
	code := createStringFromFile(routerFilePath)
	f, err := decorator.Parse(code)
	if err != nil {
		panic(err)
	}
	fn, err := getFuncByName(f, "GetRouter")
	if err != nil {
		panic(err)
	}

	// current node is a function!
	numStatements := len(fn.Body.List)
	// the last statement is a return so we insert before the return which is why it's numStatements - 1
	returnStatement := fn.Body.List[numStatements-1]
	// delete return statement
	fn.Body.List = fn.Body.List[:numStatements-1]

	controllerStmt := NewControllerStatement(data.Name)
	fn.Body.List = append(fn.Body.List, controllerStmt)

	for _, action := range data.Actions {
		routeStmt := NewRouteRegistrationStatement(action)
		fn.Body.List = append(fn.Body.List, routeStmt)
	}
	// add back return to be at the end
	fn.Body.List = append(fn.Body.List, returnStatement)

	// i don't know why i can't just directly write to the file
	// instead of using the byte buffer intermediary
	w := &bytes.Buffer{}
	if err := decorator.Fprint(w, f); err != nil {
		panic(err)
	}
	updatedContents := w.String()
	newFile, _ := os.Create(routerFilePath)
	if _, err := newFile.WriteString(updatedContents); err != nil {
		panic(err)
	}
}

func NewRouteRegistrationStatement(action Action) *dst.ExprStmt {
	return &dst.ExprStmt{
		X: &dst.CallExpr{
			Fun: &dst.SelectorExpr{
				Sel: &dst.Ident{Name: action.Method},
				X:   &dst.Ident{Name: "r"},
			},
			Args: []dst.Expr{
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s\"", action.Path),
				},
				&dst.SelectorExpr{
					Sel: &dst.Ident{Name: action.Name},
					X:   &dst.Ident{Name: fmt.Sprintf("%sCtrl", action.Resource)},
				},
			},
		},
	}
}

func NewControllerStatement(resource string) *dst.AssignStmt {
	return &dst.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []dst.Expr{
			&dst.Ident{
				Name: fmt.Sprintf("%sCtrl", resource),
			},
		},
		Rhs: []dst.Expr{
			&dst.Ident{
				Name: fmt.Sprintf("%sController{db: db, logger: log}", strings.Title(resource)),
			},
		},
	}
}

// modified from https://github.com/laher/regopher/blob/master/parsing.go
func getFuncByName(f *dst.File, funcName string) (*dst.FuncDecl, error) {
	for _, d := range f.Decls {
		switch fn := d.(type) {
		case *dst.FuncDecl:
			if fn.Name.Name == funcName {
				return fn, nil
			}
		}
	}
	return nil, fmt.Errorf("func not found")
}
