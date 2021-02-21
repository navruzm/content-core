package template

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//go:embed dist html
var embededFiles embed.FS

//Execute parses template with provided data.
func Execute(wr http.ResponseWriter, file string, layout string, data interface{}) error {
	lp, err := embededFiles.ReadFile(fmt.Sprintf("html/%s-layout.html", layout))
	if err != nil {
		log.Printf("err: %#+v\n", err)
	}

	css, err := embededFiles.ReadFile("dist/bundle.css")
	if err != nil {
		log.Printf("err: %#+v\n", err)
	}
	js, err := embededFiles.ReadFile("dist/bundle.js")
	if err != nil {
		log.Printf("err: %#+v\n", err)
	}

	lps := strings.Replace(string(lp), "{{.MainCSS}}", string(css), -1)
	lps = strings.Replace(lps, "{{.MainJS}}", string(js), -1)
	tmpl, err := template.New("layout").Funcs(tempFuncs()).Parse(lps)
	if err != nil {
		log.Printf("err: %#+v\n", err)
	}
	fp, err := embededFiles.ReadFile(fmt.Sprintf("html/%s.html", file))
	if err != nil {
		log.Printf("err: %#+v\n", err)
	}
	tmpl, err = tmpl.New(file).Parse(string(fp))
	if err != nil {
		log.Printf("err: %#+v\n", err)
	}
	// fs, err := AssetDir("html/partial")
	// if err != nil {
	// 	log.Printf("err: %#+v\n", err)
	// }
	// for _, fn := range fs {
	// 	fp, err = Asset(fmt.Sprintf("html/partial/%s", fn))
	// 	if err != nil {
	// 		log.Printf("err: %#+v\n", err)
	// 	}
	// 	tmpl, err = tmpl.New(fn).Parse(string(fp))
	// 	if err != nil {
	// 		log.Printf("err: %#+v\n", err)
	// 	}
	// }
	wr.Header().Set("Content-Type", "text/html; charset=UTF-8")
	return tmpl.ExecuteTemplate(wr, "layout", data)
}

func ExecuteString(s string, data interface{}) (string, error) {
	s = strings.ReplaceAll(s, "&quot;", "\"")
	t, err := template.New("action").Funcs(tempFuncs()).Parse(s)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func tempFuncs() template.FuncMap {
	return template.FuncMap{
		"safeHTML": func(html string) template.HTML {
			return template.HTML(html)
		},
		"safeURL": func(s string) template.URL {
			return template.URL(s)
		},
		"image": func(s string) template.HTML {
			fnp := strings.Split(s, ".")
			fp := strings.ReplaceAll(fnp[0], "-", " ")
			fp = strings.ToTitle(fp)
			var html = fmt.Sprintf(`<picture>
			    <source type="image/webp" media="(min-width: 36em)" srcset="/img/large-%s?f=webp 1024w, /img/medium-%s?f=webp 640w, /img/small-%s?f=webp 320w" sizes="33.3vw" />
			    <source type="image/webp" srcset="/img/croppedlarge-%s?f=webp 2x, /img/croppedsmall-%s?f=webp 1x" />
			    <img loading="lazy" src="/img/small-%s" alt="%s" />
			 </picture>`, s, s, s, s, s, s, fp)
			return template.HTML(html)
		},
		"thumbImage": func(s string) template.HTML {
			fnp := strings.Split(s, ".")
			fp := strings.ReplaceAll(fnp[0], "-", " ")
			fp = strings.ToTitle(fp)
			var html = fmt.Sprintf(`<picture>
          <source srcset="/img/small-%s?f=webp" type="image/webp">
          <source srcset="/img/small-%s" type="image/jpeg">
			    <img loading="lazy" width="320" width="240" src="/img/small-%s" alt="%s" />
			 </picture>`, s, s, s, fp)
			return template.HTML(html)
		},
	}
}
