package redact

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/emacski/libgosu"
	"github.com/emacski/redact/template"
)

const (
	WriteBufferSize    = 1024 * 1024 // 1MB
	PreRenderDelimeter = "RDCT_ENV"
)

// PreRenderScriptEnv attempts to run a script and extract env vars after
// allowing the script to set env vars that will be extracted
func PreRenderScriptEnv(scriptPath string) (map[string]string, error) {
	var stdout, stderr bytes.Buffer
	subcmd := fmt.Sprintf("source %s && echo '%s' && env", scriptPath, PreRenderDelimeter)
	cmd := exec.Command("sh", "-c", subcmd)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, errors.New(fmt.Sprint(stderr.String(), err))
	}
	envs := strings.Split(strings.Split(stdout.String(), PreRenderDelimeter+"\n")[1], "\n")
	return environToMap(envs[:len(envs)-1]), nil
}

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
	w := bufio.NewWriterSize(f, WriteBufferSize)
	if err = RenderCfg(tplPath, engine, w); err != nil {
		return err
	}
	w.Flush()
	return nil
}

// RenderCfgFile renders a configuration to any io.Writer
func RenderCfg(tplPath, engine string, w io.Writer) error {
	vars := GetEnvInstance().ToMap()
	eng, err := template.EngineFactory(engine)
	if err != nil {
		return err
	}
	return template.New(tplPath, vars, eng).Render(w)
}

// ExecGosu executes a command with a specific userspec using libgosu
func ExecGosu(userspec string, command []string) error {
	return libgosu.Exec(userspec, command)
}
