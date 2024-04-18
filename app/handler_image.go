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
		size, _ := bimg.NewImage(imgB).Size()
		ratio := float64(size.Width) / float64(size.Height)
		options := bimg.Options{
			Width:   size.Width,
			Height:  size.Height,
			Quality: 80,
		}
		switch typ {
		case "large", "medium":
			if ratio < 1 {
				width := int(float64(480) * ratio)
				options.Height = 480
				options.Width = width
			} else {
				height := int(float64(640) / ratio)
				options.Width = 640
				options.Height = height
			}
		case "small":
			if ratio < 1 {
				width := int(float64(320) * ratio)
				options.Height = 320
				options.Width = width
			} else {
				height := int(float64(240) / ratio)
				options.Width = 240
				options.Height = height
			}
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

		if format == "webp" {
			options.Type = bimg.WEBP
			w.Header().Add("Content-Type", "image/webp")
		}

		imgB, err = bimg.NewImage(imgB).Process(options)
		if err != nil {
			log.Error(err)
			notFoundHandler(w, r)
			return
		}
	}
	w.Header().Add("Cache-Control", "max-age=31536000")
	w.Write(imgB)
}
