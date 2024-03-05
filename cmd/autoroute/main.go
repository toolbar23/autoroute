package main

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"github.com/urfave/cli/v2"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/mod/modfile"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GetModuleName(mod_abs_path string) string {
	goModBytes, err := ioutil.ReadFile(mod_abs_path + "/" + "go.mod")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	modName := modfile.ModulePath(goModBytes)

	return modName
}

func main() {

	var watch int
	app := &cli.App{
		Name:  "autorouter",
		Usage: "make an explosive entrance",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "watch",
				Usage: "watch live, not implemented yet",
				Count: &watch,
			},
		},
		Action: func(cCtx *cli.Context) error {
			fmt.Println(cCtx.String("lang"))
			fmt.Printf("Hello %q", cCtx.Args().Get(0))
			var functype string

			git_abs_path := "/home/pm/priv/award-jury"
			mod_rel_path := "backend-go"
			routes_rel_path := "routes"
			packname := "test"
			filename := "routes.go"
			customimports := []string{}
			framework := "autoroute"

			fmt.Println("HIHIHIHIHI")
			frameworkimports := []string{}
			if framework == "gin" {
				functype = "func(c *gin.Context)"
				frameworkimports = []string{"github.com/gin-gonic/gin"}
			} else if framework == "autoroute" {
				functype = "func(c *autoroute.Context)"
				frameworkimports = []string{"github.com/gin-gonic/gin", "github.com/toolbar23/autoroute"}
			} else if framework == "http" {
				functype = "func(w http.ResponseWriter, r *http.Request, s *server.Server)"
				frameworkimports = []string{"net/http"}
			}

			// overwrite

			imports := append(customimports, frameworkimports...)

			outfilename := git_abs_path + "/" + mod_rel_path + "/" + packname + "/" + filename

			f, err := os.Create(outfilename)
			if err != nil {
				return err
			}
			defer f.Close()

			w := watcher.New()
			createRouter(f, framework, git_abs_path, mod_rel_path, routes_rel_path, w, packname, imports, functype)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}

func createRouter(writer io.Writer, framework string, git_abs_path string, mod_rel_path string, routes_rel_path string, w *watcher.Watcher, packname string, imports []string, functype string) {
	//functype := "func(g *gin/Context)"
	//imports := []string{"github.com/gin-gonic/gin"}

	//
	modAbsPath := git_abs_path + "/" + mod_rel_path
	routesAbsPath := git_abs_path + "/" + mod_rel_path + "/" + routes_rel_path
	modname := GetModuleName(modAbsPath)

	r := regexp.MustCompile(".go$")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))
	if err := w.AddRecursive(routesAbsPath); err != nil {
		log.Fatalln(err)
	}

	routes := []StaticRoute{}

	for path, _ := range w.WatchedFiles() {
		newroutes, err := processFileForRoute(modAbsPath, routesAbsPath, path)
		if err != nil {
			return
		}
		routes = append(routes, newroutes...)
	}

	writeRouter(writer, framework, packname, imports, functype, modname, routes)
}

func processFileForRoute(mod_abs_path, routes_abs_path string, file_abs_path string) ([]StaticRoute, error) {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file_abs_path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
		return []StaticRoute{}, nil
	}

	conf := types.Config{Importer: importer.Default()}
	pkg, _ := conf.Check("cmd", fset, []*ast.File{f}, nil)

	postfix, found := strings.CutPrefix(file_abs_path, routes_abs_path)
	if !found {
		return []StaticRoute{}, nil
	}
	url_path := filepath.Dir(postfix)

	postfix, found = strings.CutPrefix(file_abs_path, mod_abs_path)
	if !found {
		return []StaticRoute{}, nil
	}
	import_path := filepath.Dir(postfix)

	Routes := []StaticRoute{}
	//	rx := regexp.MustCompile("/(Get|Post|Put|Patch|Delete)(Partial|)(.*)")

	reg, _ := regexp.Compile("(Get|Put|Post|Delete|Patch)(Partial|)(.*)")
	for _, decl := range f.Decls {
		switch t := decl.(type) {
		// That's a func decl !
		case *ast.FuncDecl:

			found := reg.FindAllStringSubmatch(t.Name.Name, -1)
			fmt.Println(found)
			if len(found) > 0 {
				res := StaticRoute{}
				res.Package = pkg.Name()
				res.ImportPath = import_path
				res.Method = strings.ToUpper(found[0][1])
				res.Funcname = t.Name.Name
				res.UrlPath = url_path
				res.UrlBase = url_path
				if found[0][2] == "Partial" {
					res.Partial = found[0][3]
					res.UrlPath += "___" + res.Partial
				}
				Routes = append(Routes, res)
			}
		}
	}
	return Routes, nil
}
