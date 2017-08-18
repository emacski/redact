package template

import (
	"errors"
	"io"
	"text/template"

	"github.com/cbroglie/mustache"
)

// template engine type string names
const (
	EngineTypeGo       = "go"
	EngineTypeMustache = "mustache"
)

// Engine template engine interface
type Engine interface {
	Render(tpl *Template, w io.Writer) error
}

// EngineFactory returns an engine ptr based on the string name of the engine
func EngineFactory(name string) (Engine, error) {
	switch name {
	case EngineTypeGo:
		return &GoEngine{}, nil
	case EngineTypeMustache:
		return &MustacheEngine{}, nil
	default:
		return nil, errors.New("invalid template engine: " + name)
	}
}

// GoEngine go based template engine
type GoEngine struct {
}

// Render implements the Engine interface and renders template data
// from the io.Reader stream to the io.Writer stream
func (g *GoEngine) Render(tpl *Template, w io.Writer) error {
	tplData, err := tpl.ReadAllToString()
	if err != nil {
		return err
	}
	t, err := template.New("").Parse(tplData)
	if err != nil {
		return err
	}
	err = t.Execute(w, tpl.Vars())
	if err != nil {
		return err
	}
	return nil
}

// MustacheEngine mustache template engine
type MustacheEngine struct {
}

// Render implements the Engine interface and renders template data
// from the io.Reader stream to the io.Writer stream
func (m *MustacheEngine) Render(tpl *Template, w io.Writer) error {
	tplData, err := tpl.ReadAllToString()
	if err != nil {
		return err
	}
	r, err := mustache.Render(tplData, tpl.Vars())
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(r))
	if err != nil {
		return err
	}
	return nil
}
