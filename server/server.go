package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

const startMsg = `
Abstractify
Author: Joel Poh
 2021 Joel Poh

==> Started server, listening on port %v....
==> `

var (
	tpl *template.Template
)

type Server struct {
	Port        string
	AssetPath   string
	AssetDir    string
	TmpPath     string
	TmpDir      string
	TemplateDir string
}

func NewServer() *Server { var s Server; return &s }

func (s *Server) Start() error {
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
