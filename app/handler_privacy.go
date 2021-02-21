package app

import (
	"net/http"

	"github.com/navruzm/content-core/template"
	log "github.com/sirupsen/logrus"
)

type privacyPolicyData struct {
	templateData
}

func privacyPolicyHandler(w http.ResponseWriter, r *http.Request) {
	pageData := privacyPolicyData{
		templateData: generateTemplateData(r),
	}
	if err := template.Execute(w, "privacy-policy", "default", pageData); err != nil {
		log.Errorf("err: %s\n", err)
	}
}
