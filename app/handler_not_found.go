package app

import (
	"net/http"

	"github.com/navruzm/content-core/template"
	log "github.com/sirupsen/logrus"
)

type notFoundData struct {
	templateData
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	pageData := notFoundData{
		templateData: generateTemplateData(r),
	}
	if err := template.Execute(w, "not-found", "default", pageData); err != nil {
		log.Errorf("err: %s\n", err)
	}
}
