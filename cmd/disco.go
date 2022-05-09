package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// discoCmd represents the disco command
var discoCmd = &cobra.Command{
	Use:   "disco",
	Short: "Service discovery",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("discoverer called")
	},
}

func init() {
	rootCmd.AddCommand(discoCmd)
}
