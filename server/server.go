package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

const startMsg = `
Abstractify
Author: Joel Poh
 2021 Joel Poh

==> Started server, listening on port %v....
`

var tpl *template.Template

// Server is the struct representing the port to serve on and paths to the assets.
type Server struct {
	Port        string
	AssetPath   string
	AssetDir    string
	TmpPath     string
	TmpDir      string
	TemplateDir string
}

// NewServer returns a pointer to a Server.
func NewServer() *Server { var s Server; return &s }

// Start starts the server after it has been initialised with NewServer.
func (s *Server) Start() error {
	log.Println(strings.Repeat("-", 20))
	log.Printf(startMsg, s.Port)

	s.parseTemplates()
	s.serveStatic()
	s.setRoutes()

	err := (http.ListenAndServe(":"+s.Port, nil))
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (s *Server) parseTemplates() {
	templates := filepath.Join(s.TemplateDir, "*html")
	tpl = template.Must(template.ParseGlob(templates))
}

func (s *Server) serveStatic() {
	http.Handle(s.AssetPath, http.StripPrefix(s.AssetPath, http.FileServer(http.Dir(s.AssetDir))))
	http.Handle(s.TmpPath, http.StripPrefix(s.TmpPath, http.FileServer(http.Dir(s.TmpDir))))
}

func (s *Server) setRoutes() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/download", download)
}
