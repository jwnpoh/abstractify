package main

import (
	"log"
	"os"

	"github.com/jwnpoh/abstractify/server"
)

func main() {
	s := server.NewServer()

	s.Port = os.Getenv("PORT")
	if s.Port == "" {
		s.Port = "8080"
	}
	s.TemplateDir = "html"
	s.AssetPath = "/static/"
	s.AssetDir = "static"
	s.TmpPath = "/tmp/"
	s.TmpDir = "tmp"

	log.Fatal(s.Start())
}
