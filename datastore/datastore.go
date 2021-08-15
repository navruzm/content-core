package datastore

import (
	"embed"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Datastore interface {
	Get(slug string) (*Content, error)
	ListAll() ([]*Content, error)
	ListPosts(page int) ([]*Content, error)
	TotalPosts() int

	GetImage(name string) ([]byte, error)
}

type Content struct {
	Title           string
	Content         string
	MetaTitle       string
	MetaDescription string
	Slug            string
	Image           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewContentStore(embeddedFiles embed.FS) (Datastore, error) {
	s := ContentStorer{
		embeddedFiles: embeddedFiles,
		contentMap:    make(map[string]*Content),
		pageMap:       make(map[string]*Content),
		imageMap:      make(map[string]string),
	}
	err := s.GenerateDatastore(embeddedFiles)
	if err != nil {
		return nil, err
	}
	go func() {
		ticker := time.NewTicker(time.Hour * 24)
		for range ticker.C {
			err := s.GenerateDatastore(embeddedFiles)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	return &s, nil
}

type ContentStorer struct {
	sync.RWMutex

	embeddedFiles embed.FS
	contents      []*Content
	contentMap    map[string]*Content

	pages   []*Content
	pageMap map[string]*Content

	imageMap map[string]string
}

func (s *ContentStorer) ListPosts(page int) ([]*Content, error) {
	s.RLock()
	defer s.RUnlock()
	size := 9
	skip := (page * size) - size
	total := s.TotalPosts()
	if skip > total {
		skip = total
	}
	end := skip + size
	if end > total {
		end = total
	}
	return s.contents[skip:end], nil
}

func (s *ContentStorer) TotalPosts() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.contents)
}

func (s *ContentStorer) Get(slug string) (*Content, error) {
	s.RLock()
	defer s.RUnlock()
	if c, ok := s.contentMap[slug]; ok {
		return c, nil
	}
	if c, ok := s.pageMap[slug]; ok {
		return c, nil
	}
	return nil, errors.New("not found")
}

func (s *ContentStorer) ListAll() ([]*Content, error) {
	s.RLock()
	defer s.RUnlock()
	return append(s.contents, s.pages...), nil
}

func (s *ContentStorer) GetImage(name string) ([]byte, error) {
	s.RLock()
	defer s.RUnlock()
	if c, ok := s.imageMap[name]; ok {
		return s.embeddedFiles.ReadFile(c)
	}
	return nil, errors.New("not found")
}
