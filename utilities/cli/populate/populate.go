package populate

import (
	"errors"
	"fmt"
	"log"
	//"os"

	"github.com/Smylet/symlet-backend/api/reference"
	//"github.com/Smylet/symlet-backend/api/test"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/Smylet/symlet-backend/utilities/db"
	"github.com/Smylet/symlet-backend/utilities/utils"
)

func GetDB() (*gorm.DB, error) {
	// if os.Getenv("ENV") == "test" {
	// 	fmt.Println("Using test database")
	// 	return test.DB, nil
	// }
	
	config, err := utils.LoadConfig("../../../env")
	if err != nil {
		log.Fatal(err)
	}
	database := db.GetDB(config)
	if database == nil {
		return nil, errors.New("error connecting to database")
	}
	return database, nil
}

func PopulateReference(cmd *cobra.Command, args []string, referenceModelMap map[string]reference.ReferenceModelInterface, database *gorm.DB) error {
	var models []reference.ReferenceModelInterface
	flags, err := cmd.Flags().GetStringSlice("table")
	fmt.Println(flags, len(flags))

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
			err = errors.New("error populating " + model.GetTableName() + " table")
			return err
		}
	}
	return nil
}

func PopulateData(cmd *cobra.Command, args []string, database *gorm.DB) error {
	return db.PopulateDatabase(database)
}

var PopulateCommand = &cobra.Command{
	Use:     `populate`,
	Short:   `populate the reference table or data in the database`,
	Aliases: []string{"p"},
	Example: `populate reference --table hostel_ammenities university`,
	// Long:    `populate reference tables`,
	// PreRunE: OptionsValidator(config, headers),

}

var ReferenceCommand = &cobra.Command{
	Use:     `reference [table]...`,
	Short:   `populate the reference table or tables specified`,
	Aliases: []string{"r"},
	Example: `reference --table hostel_ammenities university`,
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := GetDB()
		if err != nil {
			return err
		}
		return PopulateReference(cmd, args, reference.ReferenceModelMap, database)
	},
}

var DataCommand = &cobra.Command{
	Use:     "data",
	Short:   "populate the data in the database",
	Aliases: []string{"d"},
	Example: `data`,
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := GetDB()
		if err != nil {
			return err
		}
		return PopulateData(cmd, args, database)
	},
}

func init() {
	ReferenceCommand.Flags().StringSliceP("table", "T", nil, "-T amenities,university")
	PopulateCommand.AddCommand(ReferenceCommand)
	PopulateCommand.AddCommand(DataCommand)
	// PopulateCommand.PersistentFlags().StringP("table", "T", "", "table to populate")
}
