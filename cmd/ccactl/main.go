package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"cca/pkg/v"
)

//go:generate go install github.com/spf13/cobra/cobra@v1.1.3
func main() {
	rootCmd := &cobra.Command{
		Use:     v.ServiceName,
		Short:   "osb service cmd",
		Long:    "cca service cmd",
		Example: "OBS [cmd]",
		Args:    cobra.MinimumNArgs(1),
		Version: v.Version,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside rootCmd Run with args: %v\n", args)
		},
	}
	subCmd1 := &cobra.Command{
		Use:   "sub1",
		Short: "subcommand1",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside subCmd1 Run with args: %v\n", args)
		},
	}
	subCmd2 := &cobra.Command{
		Use:   "sub2",
		Short: "subcommand2",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside subCmd2 Run with args: %v\n", args)
		},
	}
	subCmd11 := &cobra.Command{
		Use:   "sub11",
		Short: "subcommand11",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside subCmd11 Run with args: %v\n", args)
		},
	}
	subCmd1.AddCommand(subCmd11)
	rootCmd.AddCommand(subCmd1, subCmd2)
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
