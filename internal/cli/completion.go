package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func (c *CLI) completionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for rondo.

To load completions:

  Bash:
    source <(rondo completion bash)
    # To persist, add to ~/.bashrc or ~/.bash_profile

  Zsh:
    source <(rondo completion zsh)
    # To persist, add to ~/.zshrc

  Fish:
    rondo completion fish | source
    # To persist:
    rondo completion fish > ~/.config/fish/completions/rondo.fish

  PowerShell:
    rondo completion powershell | Out-String | Invoke-Expression
`,
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()
			switch args[0] {
			case "bash":
				return root.GenBashCompletion(os.Stdout)
			case "zsh":
				return root.GenZshCompletion(os.Stdout)
			case "fish":
				return root.GenFishCompletion(os.Stdout, true)
			case "powershell":
				return root.GenPowerShellCompletion(os.Stdout)
			default:
				return fmt.Errorf("unknown shell %q: supported shells are bash, zsh, fish, powershell", args[0])
			}
		},
	}
}
