// Date: 2021/10/22

// Package cmd
package cmd

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"caty/pkg/cmd/account"
	"caty/pkg/v"
)

func NewCmd() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:               v.ServiceName,
		Short:             "caty cli",
		Long:              "a command line tool for caty",
		SilenceErrors:     true,
		SilenceUsage:      true,
		PersistentPreRunE: initConfig,
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	persistentFlags := rootCmd.PersistentFlags()
	persistentFlags.StringP("config", "c", "", "config file (default is $HOME/.caty.yaml)")
	persistentFlags.StringP("ip", "i", "127.0.0.1:8120", "service ip (if do not provided, will lookup from config file or environment)")
	if err := viper.BindPFlag("ip", persistentFlags.Lookup("ip")); err != nil {
		return nil, err
	}
	persistentFlags.BoolP("debug", "d", false, "debug (if do not provided, will lookup from config file or environment)")
	if err := viper.BindPFlag("debug", persistentFlags.Lookup("debug")); err != nil {
		return nil, err
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Register child command
	rootCmd.AddCommand(newCompletion())
	rootCmd.AddCommand(account.NewCmd())

	return rootCmd, nil
}

func initConfig(cmd *cobra.Command, _ []string) error {
	cfg, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}
	if cfg != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfg)
	} else {
		// Find home directory.
		var home string
		if home, err = homedir.Dir(); err != nil {
			return err
		}
		// Search config in home directory with name ".woden_client" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".caty")
	}

	viper.SetEnvPrefix("cloud") // set environment variables prefix to avoid conflict
	viper.AutomaticEnv()        // read in environment variables that match

	// If a config file is found, read it in.
	if err = viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	return nil
}
