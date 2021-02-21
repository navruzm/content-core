package app

import (
	"net/http"
	"strconv"

	"github.com/navruzm/content-core/datastore"
	"github.com/navruzm/content-core/template"
	log "github.com/sirupsen/logrus"
)

type listData struct {
	templateData

	Contents   []*datastore.Content
	Content    *datastore.Content
	PageNumber int
	PrevPage   int
	NextPage   int
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageData := listData{
		templateData: generateTemplateData(r),
	}
	pageData.PageNumber, _ = strconv.Atoi(pageStr)
	if pageData.PageNumber == 0 {
		pageData.PageNumber = 1
	}
	if pageData.PageNumber > 1 {
		pageData.PrevPage = pageData.PageNumber - 1
	}
	contents, err := app.Datastore.ListPosts(pageData.PageNumber)
	if err != nil {
		log.Error(err)
	}
	pageData.Contents = contents
	if pageData.PageNumber*9 < app.Datastore.TotalPosts() {
		pageData.NextPage = pageData.PageNumber + 1
	}
	pageData.Content, _ = app.Datastore.Get("/")
	if err := template.Execute(w, "home", "default", pageData); err != nil {
		log.Errorf("err: %#+v\n", err)
	}
}
