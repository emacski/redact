package redact

import "testing"

var preRenderScriptPath = "test/pre-render.sh"

func TestPreRenderContextExec(t *testing.T) {
	ctx := new(PreRenderContext)
	envs, err := ctx.Exec(preRenderScriptPath)
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
