package app

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/h2non/bimg"
	log "github.com/sirupsen/logrus"
)

func getImage(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("f")
	image := chi.URLParam(r, "image")
	ip := strings.Split(image, "-")
	typ := ip[0]
	var err error
	imgB, err := app.Datastore.GetImage(strings.ReplaceAll(image, typ+"-", ""))
	if err != nil {
		log.Error(err)
		notFoundHandler(w, r)
		return
	}
	options := bimg.Options{
		Width:   1024,
		Height:  768,
		Quality: 70,
	}
	switch typ {
	case "medium":
		options.Width = 640
		options.Height = 480
		options.Enlarge = true
	case "small":
		options.Width = 320
		options.Height = 240
		options.Enlarge = true
	case "smallThumb":
		options.Width = 320
		options.Height = 240
		options.Crop = true
		options.Gravity = bimg.GravitySmart
	case "croppedlarge":
		options.Width = 192
		options.Height = 192
		options.Crop = true
		options.Gravity = bimg.GravitySmart
	case "croppedsmall":
		options.Width = 96
		options.Height = 96
		options.Crop = true
		options.Gravity = bimg.GravitySmart
	}

	newImage, err := bimg.NewImage(imgB).Process(options)
	if err != nil {
		log.Error(err)
		notFoundHandler(w, r)
		return
	}

	if format == "webp" {
		newWebpImage, err := bimg.NewImage(newImage).Convert(bimg.WEBP)
		if err != nil {
			log.Error(err)
		}
		if bimg.NewImage(newWebpImage).Type() == "webp" {
			newImage = newWebpImage
			w.Header().Add("Content-Type", "image/webp")
		}
	}
	w.Header().Add("Cache-Control", "max-age=31536000")
	w.Write(newImage)
}
