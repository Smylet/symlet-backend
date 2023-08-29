package root

import (
	"fmt"

	
	"github.com/Smylet/symlet-backend/utilities/cli/populate"
	"github.com/spf13/cobra"
)



var rootCmd = &cobra.Command{
	Use:     `smy`,
	Short:   `smy is the backend client for the smylet application`,
	Aliases: []string{"s"},
}

func Execute() error {
	fmt.Println("Executing root command")
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(populate.PopulateCommand)
	//rootCmd.AddCommand(migrate.MigrateCommand)
	//rootCmd.AddCommand(rollback.RollbackCommand)
	//rootCmd.AddCommand(rollback.RollbackCommand)
}