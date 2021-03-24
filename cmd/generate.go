package cmd

import (
	"github.com/spf13/cobra"
)

var useCustomMapping bool
var useCustomBinPacking bool

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates editable files and configurations needed for application",
	Long: `
Generates editable files and configurations. 

There is an optional mapping file and binpacking file that can be generated, edited and then used for the corresponding commands. If you want to use the default files, go directly to step 2.

********************************************************************************************************************************
1.a In case you want to edit the generated mapping file, first run:
  
  	  db2ctl generate mapping -c db2ctl.yaml
  
    It generates a 'mapping.csv' in 'Generated' directory, which can be edited by the user.
  
    Then, run:
  
		db2ctl generate config -c db2ctl.yaml -m 

1.b In case you want to edit the generated binpacking file, first run:
  
  	  db2ctl generate binpacking -c db2ctl.yaml
  
    It generates a 'binpacking.csv' in 'Generated' directory, which can be edited by the user.
  
    Then, run:
  
		db2ctl generate config -c db2ctl.yaml -b
********************************************************************************************************************************
  2. In case you want to use the default files, run this step:
  
  	  db2ctl generate config -c db2ctl.yaml
********************************************************************************************************************************

It generates configuration files required for pacemaker+corosync application.

`,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().StringVarP(&confFile, "conf", "c", "", "configuration yaml file needed for application")
	// generateCmd.MarkFlagRequired("conf")
	generateCmd.PersistentFlags().BoolVarP(&useCustomMapping, "map", "m", false, "(optional) use custom mapping yaml file when generating config scripts, run 'generate' first to generate mapping file")
	generateCmd.PersistentFlags().BoolVarP(&useCustomBinPacking, "bin", "b", false, "(optional) use custom bin packing file when generating config scripts, run 'generate' first to generate bin-packing file")
}
