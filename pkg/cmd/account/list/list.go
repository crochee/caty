// Date: 2021/10/22

// Package list
package list

import (
	"github.com/crochee/lib"
	"github.com/crochee/lib/table"
	"github.com/spf13/cobra"

	"cca/pkg/client"
	"cca/pkg/service/account"
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

	response, err := client.New(client.AccountService).Retrieves(cmd.Context(), opt)
	if err != nil {
		return err
	}
	listMap := make([]map[string]interface{}, len(response.Result))
	for index, value := range response.Result {
		listMap[index] = lib.Struct2Map(value)
	}
	fields := []string{
		"Verify",
		"AccountID",
		"Account",
		"UserID",
		"Email",
		"Permission",
		"Desc",
		"CreatedAt",
		"UpdatedAt",
	}
	table.RenderAsTable(listMap, fields)
	return nil
}
