package populate

import (
	"errors"
	"fmt"
	"log"

	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/spf13/cobra"

	"github.com/Smylet/symlet-backend/utilities/db"
	"github.com/Smylet/symlet-backend/utilities/utils"
)

func Populate(cmd *cobra.Command, args []string) error {
	
	referenceModelMap := map[string]reference.ReferenceModelInterface{
		"amenities": reference.ReferenceHostelAmmenities{},
		"university":        reference.ReferenceUniversity{},
	}
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
				return fmt.Errorf("invalid option %s", flag, )
			}
			models = append(models, model)
		}
	}

	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal(err)
	}
	database := db.GetDB(config)
	if database == nil {
		err = errors.New("error connecting to database")
		return err
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

func ReferenceTableArgsValidator() func(cmd *cobra.Command, args []string) error {
	referenceTableMap := map[string]reference.ReferenceModelInterface{
		"amenities": reference.ReferenceHostelAmmenities{},
		"university":        reference.ReferenceUniversity{},
	}
	return func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return nil
		}

		for _, flag := range args {
			if _, ok := referenceTableMap[flag]; !ok {
				return errors.New("invalid flag %s" + flag)
			}
		}
		return nil
	}
}


var PopulateCommand = &cobra.Command{
	Use:     `populate [table]...`,
	Short:   `populate the reference table or tables specified`,
	Aliases: []string{"p"},
	Example: `populate --table hostel_ammenities university`,
	//Long:    `populate reference tables`,
	PreRunE: ReferenceTableArgsValidator(),
	//PreRunE: OptionsValidator(config, headers),
	RunE: func(cmd *cobra.Command, args []string) error {
		return Populate(cmd, args)
	},
}



func init() {
	PopulateCommand.Flags().StringSliceP("table", "T",nil, "-T amenities,university")
	//PopulateCommand.PersistentFlags().StringP("table", "T", "", "table to populate")
}