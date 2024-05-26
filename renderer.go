package tmplrenderer

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"sync"
)

type Renderer struct {
	FS             fs.FS
	BaseDir        string
	FuncMap        template.FuncMap
	IncludePatters []string

	once      sync.Once
	mu        sync.Mutex
	templates map[string]*template.Template
}

func (r *Renderer) init() {
	r.templates = map[string]*template.Template{}
}

func (r *Renderer) Render(w io.Writer, name string, data any) error {
	r.once.Do(r.init)

	tmpl, ok := r.template(name)
	if !ok {
		tmpl, err := r.parse(name)
		if err != nil {
			return err
		}

		r.setTemplate(name, tmpl)

		return tmpl.Execute(w, data)
	}

	return tmpl.Execute(w, data)
}

func (r *Renderer) RenderBlock(w io.Writer, name, block string, data any) error {
	r.once.Do(r.init)

	tmpl, ok := r.template(name)
	if !ok {
		tmpl, err := r.parse(name)
		if err != nil {
			return err
		}

		r.setTemplate(name, tmpl)

		return tmpl.ExecuteTemplate(w, block, data)
	}

	return tmpl.ExecuteTemplate(w, block, data)
}

func (r *Renderer) parse(name string) (*template.Template, error) {
	return template.New(name).Funcs(funcMap).Funcs(r.FuncMap).ParseFS(r.FS, r.patterns(name)...)
}

func (r *Renderer) patterns(name string) []string {
	out := make([]string, len(r.IncludePatters)+1)
	for i, p := range r.IncludePatters {
		out[i] = filepath.Join(r.BaseDir, p)
	}
	out[len(out)-1] = filepath.Join(r.BaseDir, name)
	return out
}

func (r *Renderer) template(name string) (*template.Template, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.templates[name]
	return t, ok
}

func (r *Renderer) setTemplate(name string, t *template.Template) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.templates[name] = t
}

var funcMap = template.FuncMap{
	"dict": func(keyvals ...any) (map[string]any, error) {
		if len(keyvals)%2 != 0 {
			return nil, fmt.Errorf("odd number of keyvals")
		}

		out := map[string]any{}
		for i := 0; i < len(keyvals); i += 2 {
			k, ok := keyvals[i].(string)
			if !ok {
				return nil, fmt.Errorf("key not a string: %T", keyvals[i])
			}
			out[k] = keyvals[i+1]
		}

		return out, nil
	},
	"list": func(elems ...any) []any {
		return elems
	},
}
