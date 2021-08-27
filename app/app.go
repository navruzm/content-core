package app

import (
	"embed"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/navruzm/content-core/datastore"
	log "github.com/sirupsen/logrus"
)

var app *App

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(os.Stdout)
}

func NewApp(domain, sitename, analyticsID, adID string, embeddedFiles embed.FS) (*App, error) {
	var err error
	ds, err := datastore.NewContentStore(embeddedFiles)
	if err != nil {
		return nil, err
	}
	app = &App{Domain: domain, SiteName: sitename, Datastore: ds, AnalyticsID: analyticsID, AdID: adID}
	return app, nil
}

type App struct {
	Datastore   datastore.Datastore
	Domain      string
	SiteName    string
	AnalyticsID string
	AdID        string
}

func (a *App) Run(port string) error {
	if port == "" {
		port = "8080"
	}
	return http.ListenAndServe(":"+port, a.router())
}

func (a *App) router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.With(middleware.Compress(9)).HandleFunc("/", listHandler)
	r.With(middleware.Compress(9)).HandleFunc("/img/{image}", getImage)
	r.With(middleware.Compress(9)).HandleFunc("/contact", contactHandler)
	r.With(middleware.Compress(9)).HandleFunc("/terms-and-conditions", termsAndConditionsHandler)
	r.With(middleware.Compress(9)).HandleFunc("/privacy-policy", privacyPolicyHandler)
	r.With(middleware.Compress(9)).HandleFunc("/{slug}", getHandler)
	r.HandleFunc("/sitemap.xml", sitemapHandler)
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		io.Copy(w, base64.NewDecoder(base64.StdEncoding, strings.NewReader("iVBORw0KGgo=")))
	})
	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("User-agent: *\nAllow: /\nSitemap: https://%s/sitemap.xml", a.Domain)))
	})
	r.HandleFunc("/ads.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("google.com, pub-3270692583542963, DIRECT, f08c47fec0942fa0"))
	})
	r.NotFound(notFoundHandler)
	return r
}
