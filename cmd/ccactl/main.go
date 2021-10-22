package main

import (
	"os"

	"cca/pkg/cmd"
)

//go:generate go install github.com/spf13/cobra/cobra@v1.1.3
func main() {
	rootCmd, err := cmd.NewCmd()
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	if err = rootCmd.Execute(); err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}
