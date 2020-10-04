package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// generateAllConfigCmd represents the configuration command
var generateAllConfigCmd = &cobra.Command{
	Use:   "all",
	Short: "Generates all configuration files for pacemaker+corosync",
	Long: `
Generates configuration files required for pacemaker+corosync application. 

********************************************************************************************************************************
1.a In case you want to run with a custom mapping file, first run:
  
  	 db2ctl generate mapping -c db2ctl.yaml
  
    It generates a 'mapping.csv' in 'Generated' directory, which can be edited by the user.
  
    Then, run:
  
		db2ctl generate config -c db2ctl.yaml -m 

1.b In case you want to run with a custom binpacking file, first run:
  
  	 db2ctl generate binpacking -c db2ctl.yaml
  
    It generates a 'binpacking.csv' in 'Generated' directory, which can be edited by the user.
  
    Then, run:
  
		db2ctl generate config -c db2ctl.yaml -b
********************************************************************************************************************************
  2. In case you want to use the default files, run this step:
  
  	 db2ctl generate config -c db2ctl.yaml
********************************************************************************************************************************`,

	RunE: func(cmd *cobra.Command, args []string) error {
		//step 1 - parse config yaml
		//step 2 - binpacking csv - based on flag
		//step 3 - read binpacking csv
		//step 4 - genrate mapping csv
		//step 5 - read mapping csv
		//step 6 - generate configuration files
		err := command.New(cmd.Flags()).
			ParseYaml(confFile).
			GenerateBinPackingCSV().
			ReadBinPackingCSV().
			GenerateMappingFile().
			ReadMappingCSV().
			GenerateConfigFilesFromDir(""). //empty string donates generating all
			TimeTaken().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	generateCmd.AddCommand(generateAllConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateAllConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateAllConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
