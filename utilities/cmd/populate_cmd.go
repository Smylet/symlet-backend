package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/utilities/db"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var PopulateCmd = &cobra.Command{
	Use:     "populate",
	Short:   "Populate the reference table or data in the database",
	Example: `populate reference --table hostel_ammenities university`,
	RunE:    populateCmd,
}

func populateCmd(cmd *cobra.Command, args []string) error {
	database, err := initDB()
	if err != nil {
		return err
	}
	defer database.Close()

	switch cmd.Use {
	case args[0]:
		return populateReference(cmd, args, reference.ReferenceModelMap, database.GormDB())
	case args[1]:
		return populateData(cmd, args, database.GormDB())
	default:
		return errors.New("invalid command")
	}
}

func initDB() (db.DBProvider, error) {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := db.GetDB(config)
	if err != nil {
		if db != nil {
			db.Close()
		}
		return nil, err
	}

	return db, nil
}

func populateReference(cmd *cobra.Command, args []string, referenceModelMap map[string]reference.ReferenceModelInterface, database *gorm.DB) error {
	var models []reference.ReferenceModelInterface
	flags, err := cmd.Flags().GetStringSlice("table")

	if err != nil {
		return err
	}

	if len(flags) == 0 {
		for _, model := range referenceModelMap {
			models = append(models, model)
		}
	} else {
		for _, flag := range flags {
			model, ok := referenceModelMap[flag]
			if !ok {
				return fmt.Errorf("invalid option %s", flag)
			}
			models = append(models, model)
		}
	}

	for _, model := range models {
		err := model.Populate(database)
		if err != nil {
			return fmt.Errorf("error populating %s table: %w", model.GetTableName(), err)
		}
	}
	return nil
}

func populateData(cmd *cobra.Command, args []string, database *gorm.DB) error {
	return db.PopulateDatabase(database)
}
