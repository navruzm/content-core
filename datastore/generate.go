package datastore

import (
	"bytes"
	"embed"
	"sort"
	"time"

	"github.com/gosimple/slug"
	"github.com/navruzm/content-core/template"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func (s *ContentStorer) GenerateDatastore(embeddedFiles embed.FS) error {
	s.Lock()
	defer s.Unlock()
	s.contents = nil
	tn := time.Now()
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			// html.WithHardWraps(),
			// html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	files, err := embeddedFiles.ReadDir("posts")
	if err != nil {
		return err
	}
	for _, f := range files {
		content, err := embeddedFiles.ReadFile("posts/" + f.Name() + "/data.md")
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		context := parser.NewContext()
		if err := markdown.Convert(content, &buf, parser.WithContext(context)); err != nil {
			return err
		}
		t, err := template.ExecuteString(buf.String(), nil)
		if err != nil {
			return err
		}
		metaData := meta.Get(context)
		c := &Content{
			Title:           getMetaAsString("Title", metaData),
			MetaTitle:       getMetaAsString("MetaTitle", metaData),
			MetaDescription: getMetaAsString("MetaDescription", metaData),
			CreatedAt:       getMetaAsTime("CreatedAt", metaData),
			UpdatedAt:       getMetaAsTime("UpdatedAt", metaData),
			Content:         t,
			Slug:            getMetaAsString("Slug", metaData),
		}
		if tn.Before(c.CreatedAt) {
			continue
		}
		if c.Slug == "" {
			c.Slug = slug.MakeLang(c.Title, "en")
		}

		//image
		files, err = embeddedFiles.ReadDir("posts/" + f.Name())
		if err != nil {
			return err
		}
		for _, fnm := range files {
			fn := fnm.Name()
			if fn != "data.md" {
				s.imageMap[fn] = "posts/" + f.Name() + "/" + fn
				if c.Image == "" {
					c.Image = fn
				}
			}
		}

		s.contents = append(s.contents, c)
		s.contentMap[c.Slug] = c
	}

	files, err = embeddedFiles.ReadDir("pages")
	if err != nil {
		return err
	}
	for _, f := range files {
		content, err := embeddedFiles.ReadFile("pages/" + f.Name() + "/data.md")
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		context := parser.NewContext()
		if err := markdown.Convert(content, &buf, parser.WithContext(context)); err != nil {
			return err
		}
		t, err := template.ExecuteString(buf.String(), nil)
		if err != nil {
			return err
		}
		metaData := meta.Get(context)
		c := &Content{
			Title:           getMetaAsString("Title", metaData),
			MetaTitle:       getMetaAsString("MetaTitle", metaData),
			MetaDescription: getMetaAsString("MetaDescription", metaData),
			CreatedAt:       getMetaAsTime("CreatedAt", metaData),
			UpdatedAt:       getMetaAsTime("UpdatedAt", metaData),
			Content:         t,
			Slug:            getMetaAsString("Slug", metaData),
		}
		if tn.Before(c.CreatedAt) {
			continue
		}
		if c.Slug == "" {
			c.Slug = slug.MakeLang(c.Title, "en")
		}
		s.pages = append(s.pages, c)
		s.pageMap[c.Slug] = c

		//image
		files, err = embeddedFiles.ReadDir("pages/" + f.Name())
		if err != nil {
			return err
		}
		for _, fn := range files {
			if fn.Name() != "data.md" {
				s.imageMap[fn.Name()] = "pages/" + f.Name() + "/" + fn.Name()
			}
		}
	}
	sort.Slice(s.contents, func(i, j int) bool {
		return s.contents[i].UpdatedAt.After(s.contents[j].UpdatedAt)
	})
	return nil
}

func getMetaAsString(key string, metaData map[string]interface{}) string {
	if s, ok := metaData[key]; ok {
		if ss, ok := s.(string); ok {
			return ss
		}
	}
	return ""
}

func getMetaAsTime(key string, metaData map[string]interface{}) time.Time {
	if s, ok := metaData[key]; ok {
		if ss, ok := s.(string); ok {
			if ca, err := time.Parse("2006-01-02", ss); err == nil {
				return ca
			}
		}
	}
	return time.Time{}
}
