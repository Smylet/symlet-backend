package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// const (
// 	envPrefix = "SMYLET"
// )

var RootCmd = &cobra.Command{
	Use:               filepath.Base(os.Args[0]),
	Short:             "smy is the backend client for the smylet application",
	Aliases:           []string{"s"},
	PersistentPreRunE: initCmd,
	SilenceUsage:      true,
	SilenceErrors:     true,
}

func initCmd(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	level, err := log.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		return fmt.Errorf(`invalid log level "%s"`, viper.GetString("log-level"))
	}
	log.SetLevel(level)
	if log.IsLevelEnabled(log.DebugLevel) {
		log.SetReportCaller(true)
	}

	// Using common.NewLogger() if required in your case
	// logger := common.NewLogger()
	// logger.Print("Executing root command")

	return nil
}

func init() {
	RootCmd.PersistentFlags().StringP("log-level", "l", "info", "Log level")

	viper.SetConfigType("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Adding your PopulateCmd here:
	RootCmd.AddCommand(PopulateCmd)
	PopulateCmd.Flags().StringSliceP("table", "T", nil, "-T amenities,university")

	// Uncomment and add other commands as required
	// RootCmd.AddCommand(migrate.MigrateCommand)
	// RootCmd.AddCommand(rollback.RollbackCommand)
	// RootCmd.AddCommand(rollback.RollbackCommand)
}
