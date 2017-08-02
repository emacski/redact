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
func RenderCfgStdOut(tplPath string) error {
	return RenderCfg(tplPath, os.Stdout)
}

// RenderCfgFile renders a configuration to a file using the service config
func RenderCfgFile(tplPath, cfgPath string) error {
	f, err := os.Create(cfgPath)
	if err != nil {
		return err
	}
	defer f.Close()
	// use buffered writer to mitigate partially
	// written files on rendering error
	w := bufio.NewWriterSize(f, WriteBufferSize)
	defer w.Flush()
	return RenderCfg(tplPath, w)
}

// RenderCfgFile renders a configuration to any io.Writer
func RenderCfg(tplPath string, w io.Writer) error {
	return template.New(tplPath, GetEnvInstance().ToMap()).Render(w)
}

// ExecGosu executes a command with a specific userspec using libgosu
func ExecGosu(userspec string, command []string) error {
	return libgosu.Exec(userspec, command)
}
