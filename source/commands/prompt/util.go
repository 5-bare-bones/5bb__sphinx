package prompt

import (
	"errors"

	"github.com/5-bare-bones/5bb__sphinx/sig"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const template = `{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}`

// maxVisible caps the number of options rendered at once so that long lists
// scroll instead of flooding the screen.
const maxVisible = 20

// theme keeps the previous red focus / cyan help accents used with survey.
var theme = func() *huh.Theme {
	t := huh.ThemeBase()
	red := lipgloss.Color("9")
	cyan := lipgloss.Color("14")

	f := &t.Focused
	f.Title = f.Title.Foreground(red)
	f.SelectSelector = f.SelectSelector.Foreground(red)
	f.SelectedOption = f.SelectedOption.Foreground(red)
	f.SelectedPrefix = f.SelectedPrefix.Foreground(red)
	f.MultiSelectSelector = f.MultiSelectSelector.Foreground(red)
	f.Description = f.Description.Foreground(cyan)

	return t
}()

// run executes a single-field form and interrupts the process on abort.
func run(field huh.Field) error {
	form := huh.NewForm(huh.NewGroup(field)).WithTheme(theme)
	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			sig.Signal.Kill()
		}
		return err
	}

	return nil
}

// selectOne shows a single-select prompt and returns the chosen option.
func selectOne(message, help string, options []string) (string, error) {
	var choice string

	field := huh.NewSelect[string]().
		Title(message).
		Options(huh.NewOptions(options...)...).
		Value(&choice)
	if help != "" {
		field = field.Description(help)
	}
	if len(options) > maxVisible {
		field = field.Height(maxVisible + 2)
	}

	if err := run(field); err != nil {
		return "", err
	}

	return choice, nil
}

// multiSelect shows a multi-select prompt and returns the chosen options.
func multiSelect(message string, options []string) ([]string, error) {
	var choices []string

	field := huh.NewMultiSelect[string]().
		Title(message).
		Options(huh.NewOptions(options...)...).
		Value(&choices)
	if len(options) > maxVisible {
		field = field.Height(maxVisible + 2)
	}

	if err := run(field); err != nil {
		return nil, err
	}

	return choices, nil
}

// input shows a free-text prompt and returns the entered value.
func input(message, help string) (string, error) {
	var value string

	field := huh.NewInput().
		Title(message).
		Value(&value)
	if help != "" {
		field = field.Description(help)
	}

	if err := run(field); err != nil {
		return "", err
	}

	return value, nil
}
