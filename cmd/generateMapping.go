package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// mappingCmd represents the generate command
var mappingCmd = &cobra.Command{
	Use:   "mapping",
	Short: "(Optional) Generates a mapping file which can be edited before creating configuration files",
	Long: `
Generates a mapping file which can be edited before creating configuration files. 

For example:

	db2ctl generate mapping -c db2ctl.yaml
	
It generates a 'mapping.csv' in 'Generated' directory, which can be edited by the user. It can then be passed onto the 'config' command.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := command.New(cmd.Flags()).
			ParseYaml(confFile).
			GenerateMappingFile().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	generateCmd.AddCommand(mappingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mappingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
