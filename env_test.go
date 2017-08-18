package redact

import (
	"os"
	"testing"
)

func init() {
	// Set some env vars to test against
	os.Setenv("RDCT_DEFAULT_TPL_ENGINE", "mustache")
	os.Setenv("RDCT_DEFAULT_TPL_PATH", "/path/to/template")
	os.Setenv("RDCT_DEFAULT_CFG_PATH", "/path/to/config")
	os.Setenv("test_app_var", "test")
}

func TestEnvironToMap(t *testing.T) {
	envs := environToMap(os.Environ())
	if _, ok := envs["RDCT_DEFAULT_TPL_PATH"]; !ok {
		t.Error("could not find RDCT_DEFAULT_TPL_PATH in env map")
	}
	if _, ok := envs["RDCT_DEFAULT_CFG_PATH"]; !ok {
		t.Error("could not find RDCT_DEFAULT_CFG_PATH in env map")
	}
	if _, ok := envs["test_app_var"]; !ok {
		t.Error("could not find test_app_var in env map")
	}
}

func TestEnvFind(t *testing.T) {
	val := GetEnvInstance().Find("test_app_var")
	if val != "test" {
		t.Error("Expected value to be \"test\", got: ", val)
	}
	val = GetEnvInstance().Find("test_doesnt_exist")
	if len(val) > 0 {
		t.Error("Expected value to be empty, got: ", val)
	}
}

func TestEnvFindE(t *testing.T) {
	val, err := GetEnvInstance().FindE("test_app_var")
	if err != nil {
		t.Error(err)
	}
	if val != "test" {
		t.Error("Expected value to be \"test\", got: ", val)
	}
	val, err = GetEnvInstance().FindE("test_doesnt_exist")
	if err == nil {
		t.Error("Expected err to be error, got: nil")
	}
	if len(val) > 0 {
		t.Error("Expected value to be empty, got: ", val)
	}
}

func TestMerge(t *testing.T) {
	merge := map[string]string{"test_app_var": "override"}
	GetEnvInstance().Merge(merge)
	envs := GetEnvInstance().ToMap()
	if envs["test_app_var"] != "override" {
		t.Error("expected \"test_app_var\" value to be \"override\", got ", envs["test_app_var"])
	}
}

func TestEnvToMap(t *testing.T) {
	envs := GetEnvInstance().ToMap()
	if _, ok := envs["RDCT_DEFAULT_TPL_PATH"]; !ok {
		t.Error("could not find RDCT_DEFAULT_TPL_PATH in env map")
	}
	if _, ok := envs["RDCT_DEFAULT_CFG_PATH"]; !ok {
		t.Error("could not find RDCT_DEFAULT_CFG_PATH in env map")
	}
	if _, ok := envs["test_app_var"]; !ok {
		t.Error("could not find test_app_var in env map")
	}
}

func TestEnvToMapFilterPrefix(t *testing.T) {
	envs := GetEnvInstance().ToMapFilterPrefix()
	if _, ok := envs["RDCT_DEFAULT_TPL_PATH"]; !ok {
		t.Error("could not find RDCT_DEFAULT_TPL_PATH in env map")
	}
	if _, ok := envs["RDCT_DEFAULT_CFG_PATH"]; !ok {
		t.Error("could not find RDCT_DEFAULT_CFG_PATH in env map")
	}
	if len(envs) != 3 {
		t.Error("expected only 3 vars, got: ", len(envs))
	}
}

func TestEnvResolveTplEngine(t *testing.T) {
	envInstance = nil
	eng := GetEnvInstance().ResolveTplEngine()
	if eng != "mustache" {
		t.Error("expected engine to be mustache, got: ", eng)
	}
	os.Setenv("RDCT_TPL_ENGINE", "go")
	defer os.Unsetenv("RDCT_TPL_ENGINE")
	envInstance = nil
	eng = GetEnvInstance().ResolveTplEngine()
	if eng != "go" {
		t.Error("expected engine to be go, got: ", eng)
	}
}

func TestEnvResolveTplEngineDefault(t *testing.T) {
	envInstance = nil
	eng := GetEnvInstance().ResolveTplEngineDefault("go")
	if eng != "go" {
		t.Error("expected path to be go, got: ", eng)
	}
	os.Setenv("RDCT_TPL_ENGINE", "mustache")
	defer os.Unsetenv("RDCT_TPL_ENGINE")
	envInstance = nil
	eng = GetEnvInstance().ResolveCfgPathDefault("go")
	if eng != "go" {
		t.Error("expected path to be go, got: ", eng)
	}
}

func TestEnvResolveTplPath(t *testing.T) {
	envInstance = nil
	path := GetEnvInstance().ResolveTplPath()
	if path != "/path/to/template" {
		t.Error("expected path to be /path/to/template, got: ", path)
	}
	os.Setenv("RDCT_TPL_PATH", "/path/to/override/template")
	defer os.Unsetenv("RDCT_TPL_PATH")
	envInstance = nil
	path = GetEnvInstance().ResolveTplPath()
	if path != "/path/to/override/template" {
		t.Error("expected path to be /path/to/override/template, got: ", path)
	}
}

func TestResolveTplPathDefault(t *testing.T) {
	envInstance = nil
	path := GetEnvInstance().ResolveTplPathDefault("/path/to/template")
	if path != "/path/to/template" {
		t.Error("expected path to be /path/to/template, got: ", path)
	}
	os.Setenv("RDCT_TPL_PATH", "/path/to/override/template")
	defer os.Unsetenv("RDCT_TPL_PATH")
	envInstance = nil
	path = GetEnvInstance().ResolveTplPathDefault("/path/to/template")
	if path != "/path/to/override/template" {
		t.Error("expected path to be /path/to/override/template, got: ", path)
	}
}

func TestEnvResolveCfgPath(t *testing.T) {
	envInstance = nil
	path := GetEnvInstance().ResolveCfgPath()
	if path != "/path/to/config" {
		t.Error("expected path to be /path/to/config, got: ", path)
	}
	os.Setenv("RDCT_CFG_PATH", "/path/to/override/config")
	defer os.Unsetenv("RDCT_CFG_PATH")
	envInstance = nil
	path = GetEnvInstance().ResolveCfgPath()
	if path != "/path/to/override/config" {
		t.Error("expected path to be /path/to/override/config, got: ", path)
	}
}

func TestResolveCfgPathDefault(t *testing.T) {
	envInstance = nil
	path := GetEnvInstance().ResolveCfgPathDefault("/path/to/config")
	if path != "/path/to/config" {
		t.Error("expected path to be /path/to/config, got: ", path)
	}
	os.Setenv("RDCT_CFG_PATH", "/path/to/override/config")
	defer os.Unsetenv("RDCT_CFG_PATH")
	envInstance = nil
	path = GetEnvInstance().ResolveCfgPathDefault("/path/to/config")
	if path != "/path/to/override/config" {
		t.Error("expected path to be /path/to/override/config, got: ", path)
	}
}
