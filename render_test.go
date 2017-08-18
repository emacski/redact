package redact

import (
	"bytes"
	"os"
	"testing"
)

var (
	tplPathGo       = "test/test.redacted"
	tplPathMustache = "test/test.mustache"
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
