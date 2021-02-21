package app

import (
	"net/http"
)

func generateTemplateData(r *http.Request) (data templateData) {
	data = templateData{app}
	return
}

type templateData struct {
	App *App
}
