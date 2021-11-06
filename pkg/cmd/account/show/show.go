// Date: 2021/10/22

// Package show
package show

import (
	"github.com/crochee/lirity"
	"github.com/crochee/lirity/log"
	"github.com/crochee/lirity/table"
	"github.com/spf13/cobra"

	"caty/pkg/client"
	"caty/pkg/service/account"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <id>",
		Short: "Show detail of account",
		Args:  cobra.MinimumNArgs(1),
		RunE:  do,
	}

	return cmd
}

func do(cmd *cobra.Command, args []string) error {
	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return err
	}
	ctx := cmd.Context()
	if debug {
		ctx = log.WithContext(ctx, log.NewLogger(func(option *log.Option) {
			option.Level = log.DEBUG
		}))
	}
	var detail *account.RetrieveResponse
	if detail, err = client.New(client.AccountService).Retrieve(ctx, &account.User{ID: args[0]}); err != nil {
		return err
	}
	struct2Map := lirity.Struct2MapTag(detail, "")
	fields := []string{
		"UserID",
		"AccountID",
		"Account",
		"Verify",
		"Email",
		"Permission",
		"Desc",
		"CreatedAt",
		"UpdatedAt",
	}
	table.RenderAsTable(struct2Map, fields)
	return nil
}
