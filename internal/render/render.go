package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"sync"
)

type Registry struct {
	mu    sync.RWMutex
	cache map[string]*entry
}

type entry struct {
	base     *template.Template
	parseErr error
}

type Renderer struct {
	base     *template.Template
	parseErr error
}

func NewRegistry() *Registry {
	return &Registry{
		cache: make(map[string]*entry),
	}
}

func (r *Registry) LoadFS(fsys fs.FS, patterns ...string) *Renderer {
	key := fmt.Sprintf("%v%v", fsys, patterns)

	r.mu.RLock()
	cached := r.cache[key]
	r.mu.RUnlock()

	if cached != nil {
		return &Renderer{base: cached.base, parseErr: cached.parseErr}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if cached := r.cache[key]; cached != nil {
		return &Renderer{base: cached.base, parseErr: cached.parseErr}
	}

	firstFilename := ""
	if len(patterns) > 0 {
		if list, _ := fs.Glob(fsys, patterns[0]); len(list) > 0 {
			firstFilename = filepath.Base(list[0])
		}
	}

	base, err := template.New(firstFilename).
		Funcs(placeholderFuncs()).
		ParseFS(fsys, patterns...)

	cached = &entry{base: base, parseErr: err}
	r.cache[key] = cached

	return &Renderer{base: cached.base, parseErr: cached.parseErr}
}

func (r *Renderer) Render(data any, funcs template.FuncMap) (string, error) {
	if r.parseErr != nil {
		return "", r.parseErr
	}
	if r.base == nil {
		return "", errors.New("invalid template")
	}

	tmpl, err := r.base.Clone()
	if err != nil {
		return "", err
	}

	tmpl = tmpl.Funcs(funcs)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func placeholderFuncs() template.FuncMap {
	return template.FuncMap{
		"T": func(key string, args ...any) string {
			return key
		},
	}
}
