package main

import (
	"embed"
	"log"
	"to-do/app"

	_ "to-do/migrations"
)

//go:embed templates
var templates embed.FS

func main() {
	a := app.New(templates)

	if err := a.Start(); err != nil {
		log.Fatal(err)
	}
}
