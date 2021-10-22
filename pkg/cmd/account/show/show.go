// Date: 2021/10/22

// Package show
package show

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <id>",
		Short: "Show detail of area",
		Args:  cobra.MinimumNArgs(1),
		RunE:  do,
	}

	return cmd
}

func do(cmd *cobra.Command, args []string) error {
	fmt.Println(args)
	return nil
}
