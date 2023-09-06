package root

import (
	"github.com/Smylet/symlet-backend/utilities/cli/populate"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var logger = common.NewLogger()

var rootCmd = &cobra.Command{
	Use:     `smy`,
	Short:   `smy is the backend client for the smylet application`,
	Aliases: []string{"s"},
}

func Execute() error {
	logger.Print(zerolog.InfoLevel,"Executing root command")
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(populate.PopulateCommand)
	// rootCmd.AddCommand(migrate.MigrateCommand)
	// rootCmd.AddCommand(rollback.RollbackCommand)
	// rootCmd.AddCommand(rollback.RollbackCommand)
}
