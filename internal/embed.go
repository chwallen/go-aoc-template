package internal

import (
	"embed"
	"text/template"
)

var (
	//go:embed template/*
	templateFiles embed.FS
	Templates     = template.Must(template.ParseFS(templateFiles, "template/*.tmpl"))
)
