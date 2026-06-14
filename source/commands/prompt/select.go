package prompt

import (
	"fmt"
	"strings"

	"github.com/5-bare-bones/5bb__sphinx/vault/card"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"
	"github.com/5-bare-bones/5bb__sphinx/vault/file"
	"github.com/5-bare-bones/5bb__sphinx/vault/totp"
	bolt "go.etcd.io/bbolt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// selectCommands recursively looks for commands and returns a slice with all of them.
func selectCommands(parent *cobra.Command) ([]string, error) {
	if !parent.HasSubCommands() {
		return nil, nil
	}

	cmdList := parent.Commands()
	list := make([]string, len(cmdList))
	for i, c := range cmdList {
		list[i] = c.Name()
	}

	// Only when it isn't root, prepend "self",
	// which is used for executing the current command
	if parent.HasParent() {
		list = append([]string{"self"}, list...)
	}

	name, err := selectOne("Choose a command:", "", list)
	if err != nil {
		return nil, err
	}

	if name == "self" {
		return []string{}, nil
	}

	current, _, err := parent.Find([]string{name})
	if err != nil {
		return nil, err
	}

	// Repeat the process with potential childs
	child, err := selectCommands(current)
	if err != nil {
		return nil, err
	}

	result := append([]string{name}, child...)
	return result, nil
}

func selectFlags(root *cobra.Command, commands []string) ([]string, error) {
	cmd, _, err := root.Find(commands)
	if err != nil {
		return nil, errors.Wrap(err, "command not found")
	}

	if !cmd.HasFlags() {
		return nil, nil
	}

	flags, err := input("Flags:", cmd.LocalFlags().FlagUsages())
	if err != nil {
		return nil, err
	}

	return strings.Split(flags, " "), nil
}

func selectName(vault *bolt.DB, commands []string) (string, error) {
	var (
		list    []string
		err     error
		message = fmt.Sprintf("Choose a %s:", commands[0])
	)

	switch commands[0] {
	case "topt":
		list, err = totp.ListNames(vault)
		message = "Choose a TOTP:"

	case "card":
		list, err = card.ListNames(vault)

	case "file":
		list, err = file.ListNames(vault)

	case "list", "copy", "edit", "del":
		list, err = entry.ListNames(vault)
		message = "Choose an entry:"
	}
	if err != nil {
		return "", err
	}

	// If any of the commands is "list" or "topt", add the "all" option
	// to list all the elements
	for _, cmd := range commands {
		if cmd == "list" || cmd == "topt" {
			list = append([]string{"all"}, list...)
			break
		}
	}

	chosen, err := selectOne(message, "", list)
	if err != nil {
		return "", err
	}

	// The user selected all, hence return no name
	if chosen == "all" {
		chosen = ""
	}

	return chosen, nil
}

func selectManager(vault *bolt.DB) (string, error) {
	list := []string{"1Password", "Bitwarden", "Keepass", "KeepassXC", "Lastpass"}
	manager, err := selectOne("Choose a manager:", "", list)
	if err != nil {
		return "", err
	}

	return manager, nil
}

func inputName() (string, error) {
	name, err := input("Name:", "The name mustn't be empty nor include \"//\"")
	if err != nil {
		return "", err
	}

	return name, nil
}
