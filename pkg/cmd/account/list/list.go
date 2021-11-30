// Date: 2021/10/22

// Package list
package list

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
		Use:   "list",
		Short: "List account",
		RunE:  do,
	}
	cmd.Flags().StringP("account-id", "", "", "根据账户id进行搜索")
	cmd.Flags().StringP("id", "", "", "根据id进行搜索")
	cmd.Flags().StringP("account", "", "", "根据账户名进行搜索")
	cmd.Flags().StringP("email", "", "", "根据邮箱进行搜索")

	return cmd
}

func do(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	opt := &account.RetrievesRequest{}
	accountID, err := flags.GetString("account-id")
	if err != nil {
		return err
	}
	opt.AccountID = accountID
	var ID string
	if ID, err = flags.GetString("id"); err != nil {
		return err
	}
	opt.ID = ID
	var a string
	if a, err = flags.GetString("account"); err != nil {
		return err
	}
	opt.Account = a
	var email string
	if email, err = flags.GetString("email"); err != nil {
		return err
	}
	opt.Email = email

	var debug bool
	if debug, err = flags.GetBool("debug"); err != nil {
		return err
	}
	ctx := cmd.Context()
	if debug {
		ctx = log.WithContext(ctx, log.NewLogger(func(option *log.Option) {
			option.Level = log.DEBUG
		}))
	}
	var response *account.RetrieveResponses
	if response, err = client.New(client.AccountService).List(ctx, opt); err != nil {
		return err
	}
	listMap := make([]map[string]interface{}, len(response.Result))
	for index, value := range response.Result {
		listMap[index] = lirity.Struct2MapWithTag(value, "")
	}
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
	table.RenderAsTable(listMap, fields)
	return nil
}
