package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// binpackingCmd represents the binpacking command
var binpackingCmd = &cobra.Command{
	Use:   "binpacking",
	Short: "(Optional) Generates a binpacking file which can be edited before creating configuration files",
	Long: `Generates a binpacking file which can be edited before creating configuration files. 

For example:

	db2ctl generate binpacking -c db2ctl.yaml
	
It generates a 'binpacking.csv' in 'Generated' directory, which can be edited by the user. It can then be passed onto the 'config' command.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := command.New(cmd.Flags()).
			ParseYaml(confFile).
			GenerateBinPackingCSV().
			Error

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	generateCmd.AddCommand(binpackingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// binpackingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// binpackingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
