package redact

import (
	"errors"
	"os"
	"strings"
)

const (
	// prefix for reserved env vars used to configure redact itself
	EnvKeyPrefix = "RDCT_"
	// reserved env vars for redact config
	EnvKeyDefaultTplPath = "DEFAULT_TPL_PATH" // "fallback" value
	EnvKeyDefaultCfgPath = "DEFAULT_CFG_PATH" // "fallback" value
	EnvKeyTplPath        = "TPL_PATH"
	EnvKeyCfgPath        = "CFG_PATH"
)

// singleton instance
var envInstance *Env

// environToMap returns a map of all env variables
func environToMap(envs []string) map[string]string {
	envsmap := make(map[string]string)
	for _, s := range envs {
		pair := strings.Split(s, "=")
		envsmap[pair[0]] = pair[1]
	}
	return envsmap
}

// Env represents env vars as structured data
type Env struct {
	env map[string]string
}

// GetEnvInstance creates and/or returns the singleton instance of Env
func GetEnvInstance() *Env {
	if envInstance == nil {
		envInstance = &Env{environToMap(os.Environ())}
	}
	return envInstance
}

// Find returns an env var value by key
func (e *Env) Find(key string) string {
	return e.env[key]
}

// FindE returns an env var value by key or returns an error if key doesn't exist
func (e *Env) FindE(key string) (string, error) {
	val, ok := e.env[key]
	if !ok {
		return val, errors.New("key does not exist")
	}
	return val, nil
}

// ToMap returns the env vars as a map
func (e *Env) ToMap() map[string]string {
	return e.env
}

// ToMapFilterPrefix filters and returns the env vars by prefix as a map
func (e *Env) ToMapFilterPrefix() map[string]string {
	var filtered = make(map[string]string)
	for k, v := range e.env {
		if strings.HasPrefix(k, EnvKeyPrefix) {
			filtered[k] = v
		}
	}
	return filtered
}

// Merge merges key/value pairs from the supplied map to the internal env map
// where values from the supplied map will overwrite internal values given keys
// exist in both maps
func (e *Env) Merge(env map[string]string) {
	for name, val := range env {
		e.env[name] = val
	}
}

// ResolveTplPath returns the env tpl path
func (e *Env) ResolveTplPath() string {
	return e.ResolveTplPathDefault("")
}

// ResolveTplPathDefault returns the value for the template path in the
// following order: returns the env tpl path if not empty. Otherwise, returns
// the `defaultPath` value if not empty. Otherwise, returns the env default tpl
// path or empty string
func (e *Env) ResolveTplPathDefault(defaultPath string) string {
	path, err := e.FindE(EnvKeyPrefix + EnvKeyTplPath)
	if err != nil {
		if len(defaultPath) != 0 {
			return defaultPath
		}
		path = e.Find(EnvKeyPrefix + EnvKeyDefaultTplPath)
	}
	return path
}

// ResolveCfgPath returns the env cfg path
func (e *Env) ResolveCfgPath() string {
	return e.ResolveCfgPathDefault("")
}

// ResolveTplPathDefault returns the value for the config path in the
// following order: returns the env cfg path if not empty. Otherwise, returns
// the `defaultPath` value if not empty. Otherwise, returns the env default cfg
// path or empty string
func (e *Env) ResolveCfgPathDefault(defaultPath string) string {
	path, err := e.FindE(EnvKeyPrefix + EnvKeyCfgPath)
	if err != nil {
		if len(defaultPath) != 0 {
			return defaultPath
		}
		path = e.Find(EnvKeyPrefix + EnvKeyDefaultCfgPath)
	}
	return path
}
