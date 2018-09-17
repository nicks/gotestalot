package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/nicks/gotestalot/internal/server"
)

var portFlag = flag.Int("port", 8001, "Port to listen on")
var webDirFlag = flag.String("web_dir", "web", "Directory for web assets")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: gotestalot [pkg]")
	}

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	pkg := args[0]
	webDir := *webDirFlag
	s := server.NewServer(pkg, webDir)
	http.Handle("/", s.Router)

	port := *portFlag
	fmt.Printf("Starting server: http://localhost:%d/\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), http.DefaultServeMux)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server exited: %v", err)
		os.Exit(1)
	}
}
