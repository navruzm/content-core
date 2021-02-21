package main

import (
	"embed"
	"log"

	"github.com/navruzm/content-core/app"
)

//go:embed posts pages
var embededFiles embed.FS

func main() {
	a, err := app.NewApp("localhost", "Sample Site", "", "", embededFiles)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(a.Run("8080"))
}
