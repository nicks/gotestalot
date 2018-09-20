package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Server struct {
	*mux.Router

	Package string
	WebDir  string

	results map[string]*Result
	mu      *sync.Mutex
}

func NewServer(pkg, webDir string) Server {
	router := mux.NewRouter()
	s := Server{
		Package: pkg,
		WebDir:  webDir,
		Router:  router,
		mu:      &sync.Mutex{},
		results: make(map[string]*Result),
	}

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

func (s Server) createTestResult(extraArgs []string) (*Result, bool) {
	key := strings.Join(extraArgs, ",")
	s.mu.Lock()
	defer s.mu.Unlock()

	result, ok := s.results[key]
	if ok {
		return result, false
	}

	result = NewResult()
	s.results[key] = result
	return result, true
}

func (s Server) ensureTestCommandStarted(extraArgs []string) *ResultReader {
	result, isNew := s.createTestResult(extraArgs)
	if !isNew {
		return result.NewReader()
	}

	testArgs := []string{
		"test",
		"-json",
	}
	testArgs = append(testArgs, extraArgs...)
	cmd := exec.CommandContext(context.Background(), "go", testArgs...)
	stderr := bytes.NewBuffer(nil)
	cmd.Stdout = result
	cmd.Stderr = stderr

	go func() {
		err := cmd.Run()
		msg := NewSuccessMessage()
		if err != nil {
			msg = NewErrorMessage(err.Error())
			_, isExitErr := err.(*exec.ExitError)
			if isExitErr {
				msg = NewErrorMessage(fmt.Sprintf("%s\nStderr:\n%s", err, stderr.String()))
			}
		}

		encoder := json.NewEncoder(result)
		err = encoder.Encode(msg)
		if err != nil {
			log.Printf("Write: %v", err)
		}
	}()

	return result.NewReader()
}

func (s Server) streamTestCommand(res http.ResponseWriter, req *http.Request, extraArgs []string) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Printf("Upgrading websocket: %v\n", err)
		return
	}
	defer conn.Close()

	result := NewTestResultWriter(conn)
	reader := s.ensureTestCommandStarted(extraArgs)

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
}

type IndexData struct {
	Package string
}

type Message struct {
	Action string
	Output string
}

func NewErrorMessage(output string) Message {
	return Message{Action: "error", Output: output}
}

func NewSuccessMessage() Message {
	return Message{Action: "done"}
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
