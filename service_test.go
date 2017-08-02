package redact

import (
	"bytes"
	"os"
	"testing"
)

var (
	tplPath             = "test/test.redacted"
	preRenderScriptPath = "test/pre-render.sh"
)

func init() {
	os.Setenv("test_app_var", "test")
}

func TestRenderCfg(t *testing.T) {
	var rendered = new(bytes.Buffer)
	err := RenderCfg(tplPath, rendered)
	if err != nil {
		t.Error(err)
	}
	if rendered.String() != "test=test" {
		t.Error("Expected \"test=test\", got ", rendered.String())
	}
}

func TestPreRenderScriptEnv(t *testing.T) {
	envs, err := PreRenderScriptEnv(preRenderScriptPath)
	if err != nil {
		t.Error(err)
	}
	if envs["pre_render"] != "test" {
		t.Error("Expected \"pre_render\" value to be \"test\", got ", envs["pre_render"])
	}
	if envs["test_app_var"] != "override" {
		t.Error("Expected \"test_app_var\" value to be \"override\", got ", envs["test_app_var"])
	}
}
