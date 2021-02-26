package app

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/snabb/sitemap"
)

func sitemapHandler(w http.ResponseWriter, r *http.Request) {
	contents, err := app.Datastore.ListAll()
	if err != nil {
		log.Error(err)
	}
	if len(contents) == 0 {
		notFoundHandler(w, r)
		return
	}
	sm := sitemap.New()
	for _, content := range contents {
		sm.Add(&sitemap.URL{
			Loc:        strings.TrimSuffix(fmt.Sprintf("https://%s/%s", app.Domain, content.Slug), "/"),
			LastMod:    &content.UpdatedAt,
			ChangeFreq: sitemap.Monthly,
		})
	}
	sm.WriteTo(w)
}
