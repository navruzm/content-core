package app

import (
	"net/http"

	"github.com/navruzm/content-core/template"
	log "github.com/sirupsen/logrus"
)

type contactData struct {
	templateData
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	pageData := contactData{
		templateData: generateTemplateData(r),
	}
	if err := template.Execute(w, "contact", "default", pageData); err != nil {
		log.Errorf("err: %s\n", err)
	}
}
