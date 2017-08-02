package template

import (
	"io"
	"text/template"
)

// Engine template engine interface
type Engine interface {
	Render(tpl *Template, w io.Writer) error
}

// GoEngine go based template engine
type GoEngine struct {
}

// NewGoEngine constructor
func NewGoEngine() *GoEngine {
	return &GoEngine{}
}

// Render renders template data from the io.Reader stream to the
// io.Writer stream
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
