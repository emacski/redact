package redact

import (
	"bufio"
	"io"
	"os"

	"github.com/emacski/redact/template"
)

const writeBufferSize = 1024 * 1024 // 1MB

// RenderCfgStdOut renders a configuration to stdout using the service config
func RenderCfgStdOut(tplPath, engine string) error {
	return RenderCfg(tplPath, engine, os.Stdout)
}

// RenderCfgFile renders a configuration to a file using the service config
func RenderCfgFile(tplPath, cfgPath, engine string) error {
	f, err := os.Create(cfgPath)
	if err != nil {
		return err
	}
	defer f.Close()
	// use buffered writer to prevent partially written files on error
	w := bufio.NewWriterSize(f, writeBufferSize)
	if err = RenderCfg(tplPath, engine, w); err != nil {
		return err
	}
	w.Flush()
	return nil
}

// RenderCfg renders a configuration to any io.Writer
func RenderCfg(tplPath, engine string, w io.Writer) error {
	vars := GetEnvInstance().ToMap()
	eng, err := template.EngineFactory(engine)
	if err != nil {
		return err
	}
	return template.New(tplPath, vars, eng).Render(w)
}
