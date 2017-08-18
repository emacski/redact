package template

import (
	"io"
	"io/ioutil"
	"os"
)

// Template model
type Template struct {
	path   string
	vars   map[string]string
	file   *os.File
	engine Engine
}

// New creates a new template
func New(path string, vars map[string]string, engine Engine) *Template {
	return &Template{path: path, vars: vars, engine: engine}
}

// Vars getter for template vars
func (t *Template) Vars() map[string]string {
	return t.vars
}

// ReadAllToBytes reads all data from template file to byte array
func (t *Template) ReadAllToBytes() ([]byte, error) {
	if t.file == nil {
		var err error
		if t.file, err = os.Open(t.path); err != nil {
			return nil, err
		}
		defer t.file.Close()
	}
	return ioutil.ReadAll(t.file)
}

// ReadAllToString reads all data from template file to string
func (t *Template) ReadAllToString() (string, error) {
	var err error
	bytes, err := t.ReadAllToBytes()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Render renders this template to the supplied io.Writer
func (t *Template) Render(w io.Writer) error {
	return t.engine.Render(t, w)
}
