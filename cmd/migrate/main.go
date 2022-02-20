package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/crochee/lirity/command"
	"github.com/crochee/lirity/config"
	"github.com/crochee/lirity/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	rootCmd, err := NewCmd()
	if err != nil {
		log.Fatal(err)
	}
	if err = rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func NewCmd() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:               "migrate",
		Short:             "migrate cli",
		Long:              "a command line tool for migrate",
		SilenceErrors:     true,
		SilenceUsage:      true,
		PersistentPreRunE: initConfig,
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	persistentFlags := rootCmd.PersistentFlags()
	persistentFlags.StringP("config", "c", "", "config file (default is $HOME/.migrate.yaml)")
	persistentFlags.StringP("path", "p", "./pkg/migrations", "path to your migrations")
	persistentFlags.BoolP("debug", "d", false, "debug (if do not provided, will lookup from config file or environment)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Register child command
	rootCmd.AddCommand(command.NewCompletion())
	rootCmd.AddCommand(up())
	rootCmd.AddCommand(down())

	return rootCmd, nil
}

func initConfig(cmd *cobra.Command, _ []string) error {
	cfg, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}
	return config.LoadConfig(config.WithConfigFile(cfg))
}

func up() *cobra.Command {
	return &cobra.Command{
		Use: "up",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := migrateOperate(cmd, args)
			if err != nil {
				return err
			}
			defer m.Close()
			return m.Up()
		},
	}
}

func down() *cobra.Command {
	return &cobra.Command{
		Use: "down",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := migrateOperate(cmd, args)
			if err != nil {
				return err
			}
			defer m.Close()
			return m.Down()
		},
	}
}

func migrateOperate(cmd *cobra.Command, _ []string) (*migrate.Migrate, error) {
	flags := cmd.Flags()
	sourceURL, err := flags.GetString("path")
	if err != nil {
		return nil, err
	}
	o := &db.Option{
		Debug:           viper.GetString("GIN_MODE") == "debug",
		MaxOpenConn:     viper.GetInt("mysql.max_open_conns"),
		MaxIdleConn:     viper.GetInt("mysql.max_idle_conns"),
		User:            viper.GetString("mysql.user"),
		Password:        viper.GetString("mysql.password"),
		IP:              viper.GetString("mysql.ip"),
		Port:            viper.GetString("mysql.port"),
		Database:        viper.GetString("mysql.database"),
		Charset:         viper.GetString("mysql.charset"),
		ConnMaxLifetime: viper.GetDuration("mysql.conn_max_lifetime") * time.Second,
	}

	var c *sql.DB
	if c, err = sql.Open("mysql", db.Dsn(o)); err != nil {
		return nil, err
	}
	var driver database.Driver
	if driver, err = mysql.WithInstance(c, &mysql.Config{}); err != nil {
		return nil, err
	}
	return migrate.NewWithDatabaseInstance(
		"file://"+sourceURL,
		"mysql", driver)
}
