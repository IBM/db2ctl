package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// generatePacemakerCmd represents the generatePacemaker command
var generatePacemakerCmd = &cobra.Command{
	Use:   "pacemaker",
	Short: "Generate pacemaker install scripts",
	Long:  `Generate pacemaker install scripts`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := command.New(cmd.Flags()).
			ParseYaml(confFile).
			GenerateBinPackingCSV().
			ReadBinPackingCSV().
			GenerateMappingFile().
			ReadMappingCSV().
			GenerateConfigFilesFromDir(pacemakerInstallDir).
			Error

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	generateCmd.AddCommand(generatePacemakerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generatePacemakerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generatePacemakerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
