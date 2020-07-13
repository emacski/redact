# DEPRECATED

This project is **obsolete** and no longer maintained.

## ReDACT

**Reactive Docker App Configuration Toolkit**

A small (~3.0MB) command line utility with the main goal of making it simple to build (or retrofit) docker applications that are easily configured at runtime.

## Overview
ReDACT allows any docker application that uses or requires configuration files to be configured via environment variables at container runtime. This is achieved by constructing a configuration template that ReDACT will use in conjunction with environment variables to produce a config file for the application at runtime.

ReDACT is also meant to serve as the docker entrypoint for an application. This has the advantage of evaluating the environment and rendering the config template each time a container is started allowing any changes in configuration to be picked up simply by restarting a container.

# User Guide

* [Running ReDACT Containers](#running-redact-containers)
  * [Pre-Defined Configuration](#pre-defined-configuration)
  * [Custom Configuration](#custom-configuration)
* [Building ReDACT Images](#building-redact-images)
  * [Template File](#template-file)
  * [Installation](#installation)
  * [Configuration](#configuration)
  * [Pre-Render Script](#pre-render-script)
* [Example Implementations](#example-implementations)

## Running ReDACT Containers

### Pre-Defined Configuration
When running a container from an image configured with a default template, simply specifying the appropriate environment variables at container start should produce the desired configuration for the given application. Ideally, these variables should be documented with the image requiring configuration.

Example Using [k8s-kibana](https://github.com/emacski/k8s-kibana)
```bash
docker run --rm \
  -e kibana_base_url="/kibana" \
  -e kibana_elasticsearch_url="http://elasticsearch:9200" \
  emacski/kibana:latest
```

### Custom Configuration
The following environment variables are used by ReDACT to specify applicable paths at container runtime.

| Name | Description |
| ---- | ----------- |
| `RDCT_TPL_ENGINE` | Template engine to use (`go` or `mustache`). Takes precedence over `RDCT_DEFAULT_TPL_ENGINE` and cli flags. |
| `RDCT_TPL_PATH` | Path to configuration template. Takes precedence over `RDCT_DEFAULT_TPL_PATH` and cli flags. |
| `RDCT_CFG_PATH` | Location the app expects the config file. Takes precedence over `RDCT_DEFAULT_CFG_PATH` and cli flags. |

It is always possible to supply your own complete configuration at container runtime. The path to a custom configuration template can be specified by setting the `RDCT_TPL_PATH` environment variable. See [Template File](#template-file) for creating templates.

**Note:** This will often be coupled with configuring a volume that mounts the custom template to the container.

Example Using [k8s-kibana](https://github.com/emacski/k8s-kibana)
```bash
docker run --rm \
  -v $PWD/custom.yaml.redacted:/custom.yaml.redacted:ro \
  -e RDCT_TPL_PATH="/custom.yaml.redacted" \
  -e my_custom_kibana_var="some value" \
  emacski/kibana:latest
```
By default and unless otherwise specified, the default template engine used is the `go` template engine.

Example Using Custom Mustache Template
```bash
docker run --rm \
  -v $PWD/custom.yaml.mustache:/custom.yaml.mustache:ro \
  -e RDCT_TPL_ENGINE="mustache" \
  -e RDCT_TPL_PATH="/custom.yaml.mustache" \
  -e my_custom_kibana_var="some value" \
  emacski/kibana:latest
```
## Building ReDACT Images
One of the goals of ReDACT is to make the implementation as simple as possible for existing and new applications alike. In most cases, ReDACT can be implemented in the following steps:

* Create a default configuration template file
* Modify the Dockerfile to install and configure the `redact` cli utility

### Template File
The heart of ReDACT's dynamic configuration relies on config file templates which will be used to render actual application config files at container runtime. ReDACT is designed to support multiple template engines, and currently supports Go (text/template) and Mustache (https://github.com/cbroglie/mustache) templates.

See https://golang.org/pkg/text/template/ for more information on Go templates.

See https://mustache.github.io/ for more information on Mustache templates.

Given this example Go template for Kibana:
```
{{if .kibana_base_url}}
server.basePath: "{{.kibana_base_url}}"
{{end}}

{{if .kibana_elasticsearch_url}}
elasticsearch.url: "{{.kibana_elasticsearch_url}}"
{{end}}
```
ReDACT expects the environment variables `kibana_base_url` and `kibana_elasticsearch_url` to be set for modifying the config. Otherwise, the built-in Kibana defaults are used.

If `kibana_base_url` environment variable is set to `/kibana` and `kibana_elasticsearch_url` is set to `http://elasticsearch:9200`, the resulting rendered config would look like:
```
server.basePath: "/kibana"

elasticsearch.url: "http://elatsicsearch:9200"
```

When building the config template, the `redact render` command can be used to periodically check your work:
```bash
# Go template
redact render -q /path/to/template
# Or mustache template
redact render -q -e mustache /path/to/template.mustache
```
**Note:** By omitting the output flag (`-o` or `--out`) like above, the template is rendered to stdout.

### Installation
```bash
curl -L https://github.com/emacski/redact/releases/download/v0.1.0/redact -o /usr/bin/redact
chmod +x /usr/bin/redact
```
Dockerfile Example (from [k8s-kibana](https://github.com/emacski/k8s-kibana))
```dockerfile
RUN ...
  # install redact
  && curl -L https://github.com/emacski/redact/releases/download/v0.1.0/redact -o /usr/bin/redact \
  && chmod +x /usr/bin/redact \
  ...
```

### Configuration
ReDACT is aware of 4 environment variables for it's internal configuration.

| Name | Stage | Description |
| ---- | ----- | ----------- |
| `RDCT_DEFAULT_TPL_ENGINE` | Build | Default template engine to use (`go` or `mustache`). If not set, relies on cli default of `go`. |
| `RDCT_DEFAULT_TPL_PATH` | Build | File path to the default configuration template. |
| `RDCT_DEFAULT_CFG_PATH` | Build | File path to the default configuration file (the location the app expects it's config file). |
| `RDCT_TPL_ENGINE` | Run | Template engine to use (`go` or `mustache`). Takes precedence over `RDCT_DEFAULT_TPL_ENGINE` and cli flags. |
| `RDCT_TPL_PATH` | Run | File path to configuration template. Intended to be set at runtime and takes precedence over `RDCT_DEFAULT_TPL_PATH` and cli flags. |
| `RDCT_CFG_PATH` | Run | File path to configuration file location. Intended to be set at runtime and takes precedence over `RDCT_DEFAULT_CFG_PATH` and cli flags. |

Dockerfile Example (from [k8s-kibana](https://github.com/emacski/k8s-kibana))
```dockerfile
...

COPY . /

...

ENTRYPOINT ["redact", "entrypoint", \
            "--default-tpl-path", "/kibana.yml.redacted", \
            "--default-cfg-path", "/kibana/config/kibana.yml", \
            "--", \
            "kibana", "/kibana/bin/kibana"]
```
Or
```dockerfile
...

COPY . /

...

ENV RDCT_DEFAULT_TPL_PATH="/kibana.yml.redacted" \
    RDCT_DEFAULT_CFG_PATH="/kibana/config/kibana.yml"

ENTRYPOINT ["redact", "entrypoint", "--", "kibana", "/kibana/bin/kibana"]
```

**Entrypoint Command**

The `redact entrypoint` command is used to render the config template and then execute a command with a specified user spec. The above will render the template `/kibana.yml.redacted` to `/kibana/config/kibana.yml` and then execute the command `/kibana/bin/kibana` as the user:group `kibana:kibana`.

**Note:** Any flags after the `"--"` will not be parsed as `redact` command flags, but are rather assumed to be flags for the desired command being executed.

**Note:** ReDACT's command execution has the same positive side effects as using the popular `gosu` utility. In fact, ReDACT uses the `gosu` code under the hood.

### Pre-Render Script (Experimental)
ReDACT provides an experimental feature for executing a script before template rendering. The intent is that the script would set additional environment variables for additional configuration at runtime i.e. retrieving config values from an API at runtime.

Currently this feature requires the docker image to have a shell with the `source` command. Additionally, the shell and the env command should be in the PATH as `sh` and `env`. For example, while obviously bash will work, the busyboxy ash shell should also suffice. This does not impact the interpreter used to run the pre-render script as any can be used as long as it exists in the image.

## Example Implementations
The following projects may serve as useful examples. The resulting images from these projects are intended to be run in a Kubernetes cluster.

[k8s-kibana](https://github.com/emacski/k8s-kibana) - A straight forward implementation of ReDACT and Kibana.

[k8s-elasticsearch](https://github.com/emacski/k8s-elasticsearch) - A more complex implementation of ReDACT that can configure an elasticsearch cluster. This example takes advantage of the experimental pre-render script feature.
