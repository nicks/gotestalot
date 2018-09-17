package server

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router

	Package string
	WebDir  string
}

func NewServer(pkg, webDir string) Server {
	router := mux.NewRouter()
	s := Server{Package: pkg, WebDir: webDir, Router: router}

	router.HandleFunc("/", s.Index)
	router.HandleFunc("/p/{pkg}", s.ViewPackage)
	router.HandleFunc("/p/{pkg}/{name}", s.ViewTest)
	router.HandleFunc("/api/all", s.RunAll)
	router.HandleFunc("/api/p/{pkg}", s.RunPackage)
	router.HandleFunc("/api/p/{pkg}/{name}", s.RunTest)

	cssFS := http.FileServer(http.Dir(filepath.Join(s.WebDir, "css")))
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", cssFS))

	return s
}

func (s Server) templates() (*template.Template, error) {
	return template.New("server").ParseGlob(filepath.Join(s.WebDir, "templates/*"))
}

func (s Server) Index(res http.ResponseWriter, req *http.Request) {
	t, err := s.templates()
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	err = t.ExecuteTemplate(res, "index.tpl", IndexData{Package: s.Package})
	if err != nil {
		http.Error(res, err.Error(), 500)
	}
}

func (s Server) ViewPackage(res http.ResponseWriter, req *http.Request) {
	t, err := s.templates()
	if err != nil {
		http.Error(res, err.Error(), 500)
	}

	err = t.ExecuteTemplate(res, "index.tpl", IndexData{Package: s.Package})
	if err != nil {
		http.Error(res, err.Error(), 500)
	}
}

func (s Server) ViewTest(res http.ResponseWriter, req *http.Request) {
	t, err := s.templates()
	if err != nil {
		http.Error(res, err.Error(), 500)
	}

	err = t.ExecuteTemplate(res, "index.tpl", IndexData{Package: s.Package})
	if err != nil {
		http.Error(res, err.Error(), 500)
	}
}

func (s Server) RunAll(res http.ResponseWriter, req *http.Request) {
}

func (s Server) RunPackage(res http.ResponseWriter, req *http.Request) {
}

func (s Server) RunTest(res http.ResponseWriter, req *http.Request) {
}

type IndexData struct {
	Package string
}
