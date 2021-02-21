package app

import (
	"net/http"

	"github.com/navruzm/content-core/template"
	log "github.com/sirupsen/logrus"
)

type termsAndConditionsData struct {
	templateData
}

func termsAndConditionsHandler(w http.ResponseWriter, r *http.Request) {
	pageData := termsAndConditionsData{
		templateData: generateTemplateData(r),
	}
	if err := template.Execute(w, "terms-and-conditions", "default", pageData); err != nil {
		log.Errorf("err: %s\n", err)
	}
}
