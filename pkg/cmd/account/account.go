// Date: 2021/10/22

// Package account
package account

import (
	"github.com/spf13/cobra"

	"cca/pkg/cmd/account/list"
	"cca/pkg/cmd/account/show"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Manage account",
	}

	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(show.NewCmd())
	return cmd
}
