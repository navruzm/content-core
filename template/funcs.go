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
          <source media="(max-width: 799px)" srcset="/img/small-%s?f=webp" />
          <source media="(min-width: 800px)" srcset="/img/large-%s?f=webp" />
			    <img loading="lazy" src="/img/large-%s" alt="%s" />
			 </picture>`, s, s, s, fp)
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
		"simpleImage": func(s string) template.HTML {
			fnp := strings.Split(s, ".")
			fp := strings.ReplaceAll(fnp[0], "-", " ")
			fp = strings.ToTitle(fp)
			var html = fmt.Sprintf(`<img loading="lazy" src="/img/%s" alt="%s" class="img-left" />`, s, fp)
			return template.HTML(html)
		},
	}
	for name, f := range tempFuncs {
		funcs[name] = f
	}
	return funcs
}
