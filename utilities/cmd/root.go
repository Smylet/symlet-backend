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

var RootCmd = &cobra.Command{
	Use:               filepath.Base(os.Args[0]),
	Short:             "smy is the backend client for the smylet application",
	Aliases:           []string{"s"},
	PersistentPreRunE: initCmd,
	SilenceUsage:      true,
	SilenceErrors:     true,
}

func initCmd(cmd *cobra.Command, _ []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	// Check if "log-level" flag is provided, else use a default value
	logLevel := viper.GetString("log-level")
	if logLevel == "" {
		logLevel = "info" // Default to "info" log level or whatever default you deem appropriate
	}

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf(`invalid log level "%s"`, logLevel)
	}
	log.SetLevel(level)

	if log.IsLevelEnabled(log.DebugLevel) {
		log.SetReportCaller(true)
	}

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
}
