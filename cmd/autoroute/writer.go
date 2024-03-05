package main

import (
	"io"
	"strconv"
)

type StaticRoute struct {
	UrlPath    string
	UrlBase    string
	ImportPath string
	Package    string
	Funcname   string
	Method     string
	Partial    string
}

func writeRouter(w io.Writer, framework, packname string, imports []string, functype string, module string, routes []StaticRoute) {

	w.Write([]byte(`package ` + packname + `

import "regexp"
`))
	for idx, r := range routes {
		w.Write([]byte("import mod" + strconv.Itoa(idx) + " \"" + module + r.ImportPath + "\"\n"))
	}
	for _, r := range imports {
		w.Write([]byte("import  \"" + r + "\"\n"))
	}

	w.Write([]byte(`

type Route struct {
	UrlPath      string
	UrlBase    string
	Package     string
	Funcref ` + functype + `
    Funcname string
	Partial string
	Method   string
}
`))
	if framework == "http" {
		w.Write([]byte(`



	func AddRoutes(e *http.ServeMux, createHandlerFuncWithServerContext func(e Route) func(w http.ResponseWriter, r *http.Request)) {
		for _, r := range Get() {
			rx := regexp.MustCompile("/([^/]*?)--($|/)")
			path := rx.ReplaceAllString(r.UrlPath,"{$1}")
			e.HandleFunc(r.Method+" "+path, createHandlerFuncWithServerContext(r))

		}
	}
	`))

	} else {

		w.Write([]byte(`
func AddRoutes(e *gin.Engine, createHandlerFuncWithServerContext func(e Route) func(*gin.Context)) {
					rx := regexp.MustCompile("/([^/]*?)--($|/)")
	for _, r := range Get() {
			path := rx.ReplaceAllString(r.UrlPath,"/:$1")
		if r.Method == "GET" {
			e.GET(path, createHandlerFuncWithServerContext(r))
		}
		if r.Method == "POST" {
			e.POST(path, createHandlerFuncWithServerContext(r))
		}
		if r.Method == "PUT" {
			e.PUT(path, createHandlerFuncWithServerContext(r))
		}
		if r.Method == "DELETE" {
			e.DELETE(path, createHandlerFuncWithServerContext(r))
		}
	}
}
`))
	}

	w.Write([]byte(`
func Get() []Route {
	Routes := []Route{}
`))

	for idx, r := range routes {

		w.Write([]byte(`
	Routes = append(Routes, Route{
		UrlPath:     "` + r.UrlPath + `",
		UrlBase:     "` + r.UrlBase + `",
		Package:    "` + r.Package + `",
		Funcref: mod` + strconv.Itoa(idx) + "." + r.Funcname + `,
		Method:  "` + r.Method + `",
		Funcname: "` + r.Funcname + `",
		Partial: "` + r.Partial + `",
	})
`))

	}

	w.Write([]byte(`
	return Routes
}`))

}
