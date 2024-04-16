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
	var typ string
	if ip[0] == "large" || ip[0] == "medium" || ip[0] == "small" || ip[0] == "smallThumb" || ip[0] == "croppedlarge" || ip[0] == "croppedsmall" {
		typ = ip[0]
		image = strings.Replace(image, typ+"-", "", 1)
	}
	imgB, err := app.Datastore.GetImage(image)
	if err != nil {
		log.Error(err)
		notFoundHandler(w, r)
		return
	}
	if typ != "" {
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

		imgB, err = bimg.NewImage(imgB).Process(options)
		if err != nil {
			log.Error(err)
			notFoundHandler(w, r)
			return
		}
		if format == "webp" {
			imgB, err = bimg.NewImage(imgB).Convert(bimg.WEBP)
			if err != nil {
				log.Error(err)
			}
			w.Header().Add("Content-Type", "image/webp")
		}
	}
	w.Header().Add("Cache-Control", "max-age=31536000")
	w.Write(imgB)
}
