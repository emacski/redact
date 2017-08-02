package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"

	"github.com/emacski/redact"
	"github.com/spf13/cobra"
)

const help = "\n{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}"
const usageSuffix = `

{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}{{end}}
{{if .HasAvailableLocalFlags}}
Options:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} COMMAND --help" for more information about a command.{{end}}
`

// global flags
var globalQuiet bool

// render flags
var (
	renderOutPath        string
	renderScript         string
	renderDefaultTplPath string
	renderDefaultCfgPath string
)

func init() {
	rootCmd.SetHelpTemplate(help)
	rootCmd.SetUsageTemplate(usageTpl("[OPTIONS] COMMAND"))
	rootCmd.PersistentFlags().BoolVarP(&globalQuiet, "quiet", "q", false, "supress command output")

	renderCmd.SetUsageTemplate(usageTpl("[OPTIONS] [TEMPLATE_PATH]"))
	renderCmd.Flags().StringVarP(&renderOutPath, "out", "o", "", "file path to render to")
	renderCmd.Flags().StringVarP(&renderScript, "pre-render", "p", "", "EXPERIMENTAL pre-render script path")
	renderCmd.Flags().StringVar(&renderDefaultTplPath, "default-tpl-path", "", "the default template path")
	renderCmd.Flags().StringVar(&renderDefaultCfgPath, "default-cfg-path", "", "the default config path")
	rootCmd.AddCommand(renderCmd)

	execCmd.SetUsageTemplate(usageTpl("[OPTIONS] -- USERSPEC COMMAND [ARGS...]"))
	rootCmd.AddCommand(execCmd)

	entrypointCmd.SetUsageTemplate(usageTpl("[OPTIONS] -- USERSPEC COMMAND [ARGS...]"))
	entrypointCmd.Flags().StringVarP(&renderScript, "pre-render", "p", "", "EXPERIMENTAL pre-render script path")
	entrypointCmd.Flags().StringVar(&renderDefaultTplPath, "default-tpl-path", "", "the default template path")
	entrypointCmd.Flags().StringVar(&renderDefaultCfgPath, "default-cfg-path", "", "the default config path")
	rootCmd.AddCommand(entrypointCmd)

	showCmd.SetUsageTemplate(usageTpl("COMMAND"))
	rootCmd.AddCommand(showCmd)
	showEnvConfCmd.SetUsageTemplate(usageTpl(""))
	showCmd.AddCommand(showEnvConfCmd)

	versionCmd.SetUsageTemplate(usageTpl(""))
	rootCmd.AddCommand(versionCmd)
}

func usageTpl(usage string) string {
	return "Usage: {{.CommandPath}} " + usage + usageSuffix
}

func handleGlobalFlags(cmd *cobra.Command) {
	if globalQuiet {
		log.SetOutput(ioutil.Discard)
	}
}

func handlePreRenderScript(cmd *cobra.Command) error {
	if len(renderScript) != 0 {
		env := redact.GetEnvInstance()
		log.Printf(cmd.CommandPath()+": executing pre-render script %s", renderScript)
		scriptEnv, err := redact.PreRenderScriptEnv(renderScript)
		if err != nil {
			return errors.New(fmt.Sprint(cmd.CommandPath()+": ", err))
		}
		env.Merge(scriptEnv)
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:          "redact",
	Short:        "ReDACT - Reactive Docker App Configuration Toolkit",
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		handleGlobalFlags(cmd)
	},
}

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render configuration from template file",
	Long: "Render configuration from template file.\n" +
		"By default, the template is rendered to stdout",
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var env = redact.GetEnvInstance()
		// handle pre-render script
		if err = handlePreRenderScript(cmd); err != nil {
			return err
		}
		// handle config rendering
		var tplPath = env.ResolveTplPathDefault(renderDefaultTplPath)
		if len(args) != 0 {
			tplPath = args[0]
		}
		if len(tplPath) == 0 {
			return errors.New(cmd.CommandPath() + ": empty RDCT_DEFAULT_TPL_PATH or RDCT_TPL_PATH or template path arg not specified")
		}
		var cfgPath = env.ResolveCfgPathDefault(renderDefaultCfgPath)
		if len(renderOutPath) != 0 {
			cfgPath = renderOutPath
		}
		if len(cfgPath) == 0 { // no cfgPath so we render to stdout
			log.Printf(cmd.CommandPath()+": rendering template %s", tplPath)
			if err = redact.RenderCfgStdOut(tplPath); err != nil {
				return errors.New(fmt.Sprint(cmd.CommandPath()+": ", err))
			}
		} else {
			log.Printf(cmd.CommandPath()+": rendering template %s to %s", tplPath, cfgPath)
			if err = redact.RenderCfgFile(tplPath, cfgPath); err != nil {
				return errors.New(fmt.Sprint(cmd.CommandPath()+": ", err))
			}
		}
		return nil
	},
}

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command gosu style",
	Long: `USERSPEC should be: <user or uid> or <user or uid>:<group or guid>

Example: redact exec -- nobody id
         or
         redact exec -- nobody:root id`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		log.Printf(cmd.CommandPath()+": with userspec `%s`, executing command `%s`", args[0], args[1])
		if err = redact.ExecGosu(args[0], args[1:]); err != nil {
			return errors.New(fmt.Sprint(cmd.CommandPath()+": ", err))
		}
		return nil
	},
}

var entrypointCmd = &cobra.Command{
	Use:   "entrypoint",
	Short: "Renders configuration then executes a command",
	Long: `USERSPEC should be: <user or uid> or <user or uid>:<group or guid>

Example: redact entrypoint -- nobody id
         or
         redact entrypoint -- nobody:root id`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var env = redact.GetEnvInstance()
		// handle pre-render script
		if err = handlePreRenderScript(cmd); err != nil {
			return err
		}
		// handle config rendering
		var tplPath = env.ResolveTplPathDefault(renderDefaultTplPath)
		if len(tplPath) == 0 {
			return errors.New(cmd.CommandPath() + ": empty RDCT_DEFAULT_TPL_PATH or RDCT_TPL_PATH or --default-tpl-path not specified")
		}
		var cfgPath = env.ResolveCfgPathDefault(renderDefaultCfgPath)
		if len(cfgPath) == 0 {
			return errors.New(cmd.CommandPath() + ": empty RDCT_DEFAULT_CFG_PATH or RDCT_CFG_PATH or --default-cfg-path not specified")
		}
		log.Printf(cmd.CommandPath()+": rendering template %s to %s", tplPath, cfgPath)
		if err = redact.RenderCfgFile(tplPath, cfgPath); err != nil {
			return errors.New(fmt.Sprint(cmd.CommandPath()+": ", err))
		}
		// handle command execution
		log.Printf(cmd.CommandPath()+": with userspec `%s`, executing command `%s`", args[0], args[1])
		if err = redact.ExecGosu(args[0], args[1:]); err != nil {
			return errors.New(fmt.Sprint(cmd.CommandPath()+": ", err))
		}
		return nil
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Debugging and troubleshooting outputs",
}

var showEnvConfCmd = &cobra.Command{
	Use:   "env-config",
	Short: "Show redact environment config",
	Run: func(cmd *cobra.Command, args []string) {
		envs := redact.GetEnvInstance().ToMapFilterPrefix()
		format := fmt.Sprintf("%%-%ds%%s", func() (w int) {
			for k, _ := range envs {
				if len(k) > w {
					w = len(k)
				}
			}
			return w + 1 // pad one col
		}())
		for name, val := range envs {
			log.Printf(format, name, val)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("redact version %s %s %s %s", redact.Version, runtime.GOOS,
			runtime.GOARCH, runtime.Compiler)
	},
}