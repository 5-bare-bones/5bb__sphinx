package root

import (
	"regexp"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// Palette shared with the interactive prompts (red accent / cyan options) so the
// CLI looks consistent whether you are reading help or answering a prompt.
var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
	cmdStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	flagStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	faintStyle  = lipgloss.NewStyle().Faint(true)
)

// flagToken matches short (-x) and long (--name) flag names so they can be
// highlighted inside the otherwise pre-formatted FlagUsages block.
var flagToken = regexp.MustCompile(`-{1,2}[A-Za-z][\w-]*`)

func init() {
	cobra.AddTemplateFuncs(map[string]interface{}{
		"header":  headerStyle.Render,
		"cmdname": cmdStyle.Render,
		"faint":   faintStyle.Render,
		"flags": func(usage string) string {
			return flagToken.ReplaceAllStringFunc(usage, func(s string) string {
				return flagStyle.Render(s)
			})
		},
	})
}

// applyHelpTheme installs the colorized usage template on the command. cobra
// resolves a child's template by walking up to its parent, so setting it on the
// root is enough to colorize every command's --help output.
func applyHelpTheme(cmd *cobra.Command) {
	cmd.SetUsageTemplate(colorUsageTemplate)
}

// colorUsageTemplate mirrors cobra's defaultUsageTemplate, wrapping the section
// headers, command names and flag names in the styles above. Keep its structure
// in sync with cobra's default if the dependency is upgraded.
const colorUsageTemplate = `{{header "Usage:"}}{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

{{header "Aliases:"}}
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

{{header "Examples:"}}
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

{{header "Available Commands:"}}{{range $cmds}}{{if .IsAvailableCommand}}
  {{cmdname (rpad .Name .NamePadding)}} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{header $group.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) .IsAvailableCommand)}}
  {{cmdname (rpad .Name .NamePadding)}} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

{{header "Additional Commands:"}}{{range $cmds}}{{if (and (eq .GroupID "") .IsAvailableCommand)}}
  {{cmdname (rpad .Name .NamePadding)}} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

{{header "Flags:"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces | flags}}{{end}}{{if .HasAvailableInheritedFlags}}

{{header "Global Flags:"}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces | flags}}{{end}}{{if .HasHelpSubCommands}}

{{header "Additional help topics:"}}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{cmdname (rpad .CommandPath .CommandPathPadding)}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

{{faint "Use"}} "{{.CommandPath}} [command] --help" {{faint "for more information about a command."}}{{end}}
`
