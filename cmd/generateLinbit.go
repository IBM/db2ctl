package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// generateLinbitCmd represents the generateLinbit command
var generateLinbitCmd = &cobra.Command{
	Use:   "linbit",
	Short: "Generate linbit install scripts",
	Long:  `Generate linbit install scripts`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := command.New(cmd.Flags()).
			ParseYaml(confFile).
			GenerateBinPackingCSV().
			ReadBinPackingCSV().
			GenerateConfigFilesFromDir(linbitInstallDir).
			Error

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	generateCmd.AddCommand(generateLinbitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateLinbitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateLinbitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
