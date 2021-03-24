package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// installAllCmd represents the installAll command
var installAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Install all modules",
	Long: `Install all modules. For example:

	db2ctl install all -c db2ctl.yaml"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := command.New(cmd.Flags()).
			ParseYaml(confFile).
			GenerateBinPackingCSV().
			ReadBinPackingCSV().
			GenerateMappingFile().
			ReadMappingCSV().
			GenerateConfigFilesFromDir("install").
			RunBashScripts().
			DeleteStateForDir("cleanup").
			TimeTaken().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	installCmd.AddCommand(installAllCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installAllCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installAllCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
