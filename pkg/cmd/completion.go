// Date: 2021/10/22

// Package cmd
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func newCompletion() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powerfish]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

$ source <(yourprogram completion bash)

# To load completions for each session, execute once:
Linux:
  $ yourprogram completion bash > /etc/bash_completion.d/yourprogram
MacOS:
  $ yourprogram completion bash > /usr/local/etc/bash_completion.d/yourprogram

Zsh:

$ source <(yourprogram completion zsh)

# To load completions for each session, execute once:
$ yourprogram completion zsh > "${fpath[1]}/_yourprogram"
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletionV2(os.Stdout, false)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, false)
			case "powershell":
				return cmd.Root().GenPowerShellCompletion(os.Stdout)
			default:
				return nil
			}
		},
	}
}
