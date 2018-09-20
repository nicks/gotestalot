package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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

	jsFS := http.FileServer(http.Dir(filepath.Join(s.WebDir, "js")))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", jsFS))

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
	s.streamTestCommand(res, req, []string{fmt.Sprintf("%s/...", s.Package)})
}

func (s Server) RunPackage(res http.ResponseWriter, req *http.Request) {
}

func (s Server) RunTest(res http.ResponseWriter, req *http.Request) {
}

func (s Server) streamTestCommand(res http.ResponseWriter, req *http.Request, extraArgs []string) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Printf("Upgrading websocket: %v\n", err)
		return
	}
	defer conn.Close()

	result := NewTestResultWriter(conn)
	testArgs := []string{
		"test",
		"-json",
	}
	testArgs = append(testArgs, extraArgs...)

	cmd := exec.CommandContext(req.Context(), "go", testArgs...)
	reader, writer := io.Pipe()
	stderr := bytes.NewBuffer(nil)
	cmd.Stdout = writer
	cmd.Stderr = stderr
	done := make(chan struct{})

	go func() {
		defer close(done)

		decoder := json.NewDecoder(reader)
		for decoder.More() {
			var msg json.RawMessage
			err := decoder.Decode(&msg)
			if err != nil {
				result.WriteError(err.Error())
				continue
			}

			result.WriteJSON(msg)
		}
	}()

	err = cmd.Run()
	if err != nil {
		_, isExitErr := err.(*exec.ExitError)
		if isExitErr {
			result.WriteError(fmt.Sprintf("%s\nStderr:\n%s", err, stderr.String()))
		} else {
			result.WriteError(err.Error())
		}
	}

	<-done
}

type IndexData struct {
	Package string
}

type ErrorMessage struct {
	Action string
	Output string
}

func NewErrorMessage(output string) ErrorMessage {
	return ErrorMessage{Action: "error", Output: output}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type TestResultWriter struct {
	conn *websocket.Conn
	mu   *sync.Mutex
}

func NewTestResultWriter(conn *websocket.Conn) TestResultWriter {
	return TestResultWriter{
		conn: conn,
		mu:   &sync.Mutex{},
	}
}

func (w TestResultWriter) WriteJSON(msg json.RawMessage) {
	w.writeInternal(msg)
}

func (w TestResultWriter) WriteError(output string) {
	w.writeInternal(NewErrorMessage(output))
}

func (w TestResultWriter) writeInternal(obj interface{}) {
	encoded, err := json.Marshal(obj)
	if err != nil {
		log.Printf("WriteError: %v", err)
		return
	}

	w.mu.Lock()
	err = w.conn.WriteMessage(websocket.TextMessage, encoded)
	w.mu.Unlock()

	if err != nil {
		log.Printf("WriteError: %v", err)
	}
}
