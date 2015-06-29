package render

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path"
)

// Example:
//
// rdr := NewRenderer("../admin/templates", "../templates")
// rdr.AddTemplateExtension(".tmpl")
// rdr.RegisterTemplateFunc("cvtcolor", func(color rgb) string { ... })
// if err := rdr.ParseFiles(); err != nil {
//    log.Fatalf("init renderer failed, error: %v", err)
// }
// // ...
// if err := rdr.Render(rw, "home_page", data); err != nil {
//    log.Errorf("failed to render page, error: %v", err)
// }

type IRender interface {
	Render(wr io.Writer, template string, data interface{}) error
}

type Renderer struct {
	dirs            []string
	err             error
	parsed          bool
	templates       map[string]*template.Template
	templateFuncs   template.FuncMap
	templateExtList map[string]bool
}

func NewRenderer(dirs ...string) *Renderer {
	rdr := &Renderer{
		dirs:            dirs,
		templates:       make(map[string]*template.Template, 25),
		templateFuncs:   make(template.FuncMap),
		templateExtList: map[string]bool{".tpl": true, ".html": true},
	}
	return rdr
}

func (rdr *Renderer) Render(wr io.Writer, template string, data interface{}) error {
	if template == "" {
		return errors.New("empty template name")
	}

	tpl, err := rdr.Lookup(template)
	if err != nil {
		return err
	}

	return tpl.Execute(wr, data)
}

// By default, only files with ".tpl" or ".html" extension will be parsed as
// template files.
// NB: Do this before parsing the template files.
func (rdr *Renderer) AddTemplateExtension(exts ...string) {
	if rdr.parsed {
		return
	}
	for _, ext := range exts {
		rdr.templateExtList[ext] = true
	}
}

// Register template functions to renderer.
// NB: Do this before parsing the template files.
// More about `FuncMap`, see http://godoc.org/text/template#FuncMap.
func (rdr *Renderer) RegisterTemplateFunc(name string, fn interface{}) {
	if rdr.parsed {
		return
	}
	rdr.templateFuncs[name] = fn
}

// Start parsing the files under the specified template dir.
// Returns parse log and error.
func (rdr *Renderer) ParseFiles() error {
	if rdr.parsed {
		return nil
	}

	if len(rdr.dirs) == 0 {
		return errors.New("no directories set")
	}

	// Reset internal error to nil. Maybe it's a second time call.
	rdr.err = nil

	for _, dir := range rdr.dirs {
		rdr.traverse(dir)
		if rdr.err != nil {
			break
		}
	}
	if rdr.err == nil {
		rdr.parsed = true
	}

	return rdr.err
}

// Find a template by name. Returns nil if not found.
func (rdr *Renderer) Lookup(name string) (*template.Template, error) {
	if tpl, ok := rdr.templates[name]; ok {
		return tpl, nil
	} else if !rdr.parsed {
		return nil, errors.New("templates not parsed")
	} else {
		return nil, fmt.Errorf("template %q not found", name)
	}
}

// Traverse a directory, and find the target files to parse.
func (rdr *Renderer) traverse(templateDir string) {
	if rdr.err != nil {
		return
	}

	var scandir func(dir string) []string

	scandir = func(dir string) (files []string) {
		if rdr.err != nil {
			return
		}
		if fs, err := ioutil.ReadDir(dir); err != nil {
			rdr.err = fmt.Errorf("unable to read directory %q during traversing, due to %v", dir, err.Error())
		} else {
			for _, fitem := range fs {
				filename := path.Join(dir, fitem.Name())
				if fitem.IsDir() {
					files = append(files, scandir(filename)...)
				} else if rdr.templateExtList[path.Ext(filename)] {
					files = append(files, filename)
				}
			}
		}
		return
	}

	files := scandir(templateDir)
	if rdr.err != nil || len(files) == 0 {
		return
	}

	t, err := template.New("/fuck/").Funcs(rdr.templateFuncs).ParseFiles(files...)
	if err != nil {
		rdr.err = fmt.Errorf("parse template file error, due to %v", err.Error())
		return
	}
	for _, it := range t.Templates() {
		if it.Name() == "/fuck/" {
			continue
		}
		rdr.templates[it.Name()] = it
	}
}
