// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/28

package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"obs/cmd"
)

// todo build cli
func main() {
	rootCmd := &cobra.Command{
		Version: cmd.Version,
		Use:     cmd.ServiceName,
		Aliases: []string{"run"},
		Short:   "run cmd",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside rootCmd Run with args: %v\n", args)
		},
	}
	subCmd := &cobra.Command{
		Use:   "sub [no options!]",
		Short: "My subcommand",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside subCmd Run with args: %v\n", args)
		},
	}
	rootCmd.AddCommand(subCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Print(err)
	}
}
