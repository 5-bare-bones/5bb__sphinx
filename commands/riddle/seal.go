package riddle

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newSealCmd() *cobra.Command {
	var prompt, answer, secret, out string
	cmd := &cobra.Command{
		Use:   "seal",
		Short: "Seal a secret behind a riddle answer",
		Example: `
sphinx riddle seal --prompt "I speak without a mouth..." --answer "an echo" --secret <token> -o echo.riddle`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if prompt == "" || answer == "" || secret == "" {
				return errors.New("--prompt, --answer and --secret are required")
			}
			if out == "" {
				return errors.New("--out is required")
			}

			rf, err := seal(prompt, answer, []byte(secret))
			if err != nil {
				return errors.Wrap(err, "sealing riddle")
			}

			data, err := json.MarshalIndent(rf, "", "  ")
			if err != nil {
				return err
			}
			if err := os.WriteFile(out, data, 0o600); err != nil {
				return errors.Wrap(err, "writing riddle file")
			}

			fmt.Fprintf(cmd.OutOrStdout(), "sealed riddle into %s\n", out)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVar(&prompt, "prompt", "", "the riddle text shown to the solver")
	f.StringVar(&answer, "answer", "", "the answer whose normalized form derives the key")
	f.StringVar(&secret, "secret", "", "the token revealed on a correct answer")
	f.StringVarP(&out, "out", "o", "", "output riddle file")
	return cmd
}
