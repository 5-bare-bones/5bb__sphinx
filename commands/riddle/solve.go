package riddle

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/5-bare-bones/5bb__sphinx/terminal"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newSolveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "solve <riddle-file>",
		Short: "Answer a riddle to reveal its sealed token",
		Example: `
sphinx riddle solve echo.riddle`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := os.ReadFile(args[0])
			if err != nil {
				return errors.Wrap(err, "reading riddle file")
			}

			var rf riddleFile
			if err := json.Unmarshal(data, &rf); err != nil {
				return errors.Wrap(err, "parsing riddle file")
			}

			fmt.Fprintln(cmd.OutOrStdout(), rf.Prompt)
			answer := terminal.ScanOneLine(bufio.NewReader(cmd.InOrStdin()), "Answer")

			secret, err := solve(rf, answer)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "correct! token: %s\n", secret)
			return nil
		},
	}
}
