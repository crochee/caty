// Date: 2021/10/22

// Package list
package list

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List account",
		RunE:  do,
	}

	return cmd
}

func do(cmd *cobra.Command, args []string) error {
	fmt.Println(args)
	return nil
}
