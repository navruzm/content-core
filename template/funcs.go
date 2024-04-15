package template

import (
	"fmt"
	"html/template"
	"strings"
)

var tempFuncs = template.FuncMap{}

func RegisterFuncs(funcs template.FuncMap) {
	for k, v := range funcs {
		tempFuncs[k] = v
	}
}

func getTempFuncs() template.FuncMap {
	funcs := template.FuncMap{
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
          <source srcset="/img/smallThumb-%s?f=webp" type="image/webp">
          <source srcset="/img/smallThumb-%s" type="image/jpeg">
			    <img loading="lazy" src="/img/smallThumb-%s" alt="%s" />
			 </picture>`, s, s, s, fp)
			return template.HTML(html)
		},
	}
	for name, f := range tempFuncs {
		funcs[name] = f
	}
	return funcs
}
