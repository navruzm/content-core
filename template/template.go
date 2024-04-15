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

// Execute parses template with provided data.
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
	tmpl, err := template.New("layout").Funcs(getTempFuncs()).Parse(lps)
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
	t, err := template.New("action").Funcs(getTempFuncs()).Parse(s)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
