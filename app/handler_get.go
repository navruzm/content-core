package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/navruzm/content-core/datastore"
	"github.com/navruzm/content-core/template"
	log "github.com/sirupsen/logrus"
)

type getData struct {
	templateData
	Content  *datastore.Content
	Category string
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	pageData := getData{
		templateData: generateTemplateData(r),
	}
	var err error
	pageData.Content, err = app.Datastore.Get(slug)
	if err != nil {
		log.Error(err)
		notFoundHandler(w, r)
		return
	}
	if err := template.Execute(w, "content", "default", pageData); err != nil {
		log.Errorf("err: %s\n", err)
	}
}
