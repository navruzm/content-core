package app

import (
	"net/http"
)

func generateTemplateData(_ *http.Request) (data templateData) {
	data = templateData{app}
	return
}

type templateData struct {
	App *App
}
