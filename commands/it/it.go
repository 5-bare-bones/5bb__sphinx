package it

import (
	"strings"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* No arguments
sphinx it

* Command without flags
sphinx it list

* Command with flags
sphinx it list -s -q

* Only the name
sphinx sample`

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "it <command|flags|name>",
		Short: "Execute commands through an interactive prompt",
		Long: `Interactive prompt.
This command behaves depending on the arguments received, it requests the missing information.

Given 				Requests
command 			flags and name
command and flags 		name
name 				command and flags`,
		Example:            example,
		DisableFlagParsing: true,
		RunE:               runIt(vault),
	}
}

func runIt(vault *bolt.DB) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		root := cmd.Root()
		// Get rid of unnecessary information and reset in case we are inside a session
		defer root.SetUsageTemplate(root.UsageTemplate())
		root.SetUsageTemplate(template)

		// We received nothing, request all
		if len(args) == 0 {
			arguments, err := requestCommands(vault, root, nil)
			if err != nil {
				return err
			}

			return execute(root, arguments)
		}

		command, _, err := root.Find(args)
		if err != nil || command == root {
			// If the command does not exist or is the root, assume the user passed a name
			arguments, err := gotName(vault, root, args)
			if err != nil {
				return err
			}

			return execute(root, arguments)
		}

		foundFlags := false
		for _, a := range args {
			if strings.HasPrefix(a, "-") {
				foundFlags = true
				break
			}
		}

		if foundFlags {
			// Got command+flags, do not look for subcommands
			if err := command.ParseFlags(args); err != nil {
				return err
			}

			// Get rid of the command and flags to validate the name
			argsWoFlags := strings.Join(command.Flags().Args(), " ")
			name := strings.Replace(argsWoFlags, command.Name(), "", 1)

			// The validation won't fail if the user lists records
			err := command.ValidateArgs([]string{name})
			if err != nil || strings.Contains(command.Name(), "list") {
				// Received commands+flags, request name
				arguments, err := requestName(vault, args)
				if err != nil {
					return err
				}

				return execute(root, arguments)
			}

			// Received command+flags+name, nothing to request
			return execute(root, args)
		}

		// Pass on received command(s) and look for subcommands
		arguments, err := requestCommands(vault, command, args)
		if err != nil {
			return err
		}

		return execute(root, arguments)
	}
}

func execute(root *cobra.Command, args []string) error {
	// Discard empty arguments as some commands will fail if we don't
	// eg. file show
	filteredArgs := make([]string, 0, len(args))
	for _, a := range args {
		if a != "" {
			filteredArgs = append(filteredArgs, a)
		}
	}

	root.SetArgs(filteredArgs)
	return root.Execute()
}

func requestCommands(vault *bolt.DB, root *cobra.Command, receivedCmds []string) ([]string, error) {
	commands, err := selectCommands(root)
	if err != nil {
		return nil, err
	}

	flags, err := selectFlags(root, commands)
	if err != nil {
		return nil, err
	}

	args := append(commands, flags...)
	// Prepend the received commands if there is any
	// We would have [received commands] [commands] [flags]
	if len(receivedCmds) > 0 {
		args = append(receivedCmds, args...)
	}
	return requestName(vault, args)
}

// args contains commands and flags.
func requestName(vault *bolt.DB, args []string) ([]string, error) {
	var (
		name string
		err  error
	)

	search := strings.Join(args, " ")
	// contains reports whether s is within search
	contains := func(s string) bool {
		return strings.Contains(search, s)
	}

	// Behave depending on which command the user is executing
	switch {
	case contains("add"),
		contains("list") && contains("-f"), // Filter
		contains("del") && contains("-d"):  // Remove directory
		name, err = inputName()

	case contains("import"), contains("export"):
		name, err = selectManager(vault)

	case contains("file show"), contains("file touch"):
		names, err := fileMultiselect(vault)
		if err != nil {
			return nil, err
		}
		return append(args, names...), nil

	case contains("file move"):
		names, err := fileMvNames(vault)
		if err != nil {
			return nil, err
		}
		return append(args, names...), nil

	default:
		list := []string{"topt", "copy", "edit", "list", "del"}
		// Request the name depending on the command
		for _, cmd := range list {
			if contains(cmd) {
				// Skip "config edit" as it doesn't need a name
				if args[0] != "config" {
					name, err = selectName(vault, args)
					break
				}
			}
		}
	}
	if err != nil {
		return nil, err
	}

	// Remember: the flags are inside the commands slice
	result := append(args, name)
	return result, nil
}

// gotName is executed when the user already provided the name, commands and flags are requested only.
func gotName(vault *bolt.DB, root *cobra.Command, args []string) ([]string, error) {
	var (
		name  []string
		flags []string
	)

	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			flags = append(flags, a)
			continue
		}
		name = append(name, a)
	}

	commands, err := selectCommands(root)
	if err != nil {
		return nil, err
	}

	if len(flags) == 0 {
		flags, err = selectFlags(root, commands)
		if err != nil {
			return nil, err
		}
	}

	if len(name) == 0 {
		name, err = requestName(vault, commands)
		if err != nil {
			return nil, err
		}
	}

	result := append(commands, flags...)
	result = append(result, name...)

	return result, nil
}
