package redact

import (
	"bytes"
	"os"
	"testing"
)

var (
	tplPathGo           = "test/test.redacted"
	tplPathMustache     = "test/test.mustache"
	preRenderScriptPath = "test/pre-render.sh"
)

func init() {
	os.Setenv("test_app_var", "test")
}

func TestRenderCfgGoEngine(t *testing.T) {
	var rendered = new(bytes.Buffer)
	err := RenderCfg(tplPathGo, "go", rendered)
	if err != nil {
		t.Error(err)
	}
	if rendered.String() != "test=test\n" {
		t.Error("Expected \"test=test\", got: ", rendered.String())
	}
}

func TestRenderCfgMustacheEngine(t *testing.T) {
	var rendered = new(bytes.Buffer)
	err := RenderCfg(tplPathMustache, "mustache", rendered)
	if err != nil {
		t.Error(err)
	}
	if rendered.String() != "test=test\n" {
		t.Error("Expected \"test=test\", got: ", rendered.String()[9])
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
