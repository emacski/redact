package redact

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const preRenderDelimeter = "RDCT_PRERENDER_ENV"

// PreRenderContext represents the execution context for a pre-render script
type PreRenderContext struct {
	StdOut string
	StdErr string
	stdout bytes.Buffer // raw stdout from script execution
	stderr bytes.Buffer // raw stderr from script execution
}

// Exec executes the pre-render script within the context
func (p *PreRenderContext) Exec(scriptPath string) (map[string]string, error) {
	p.reset() // reset std stream states before each run
	subcmd := fmt.Sprintf("source %s && echo '%s' && env", scriptPath, preRenderDelimeter)
	cmd := exec.Command("sh", "-c", subcmd)
	cmd.Stdout, cmd.Stderr = &p.stdout, &p.stderr
	if err := cmd.Run(); err != nil {
		p.StdOut, p.StdErr = p.stdout.String(), p.stderr.String()
		return nil, errors.New(fmt.Sprint(p.stderr.String(), err))
	}
	stdoutSplit := strings.Split(p.stdout.String(), preRenderDelimeter+"\n")
	p.StdOut, p.StdErr = stdoutSplit[0], p.stderr.String()
	envlist := strings.Split(stdoutSplit[1], "\n")
	return environToMap(envlist[:len(envlist)-1]), nil
}

func (p *PreRenderContext) reset() {
	p.StdOut, p.StdErr = "", ""
	p.stdout, p.stderr = bytes.Buffer{}, bytes.Buffer{}
}
