/*
Copyright © 2022 Juanma
*/
package cmd

import (
	"fmt"

	"github.com/rocajuanma/anvil/pkg/figure"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Batch installation of applications using brew",
	Long:  SETUP_COMMAND_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setup called")
	},
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push assets to Github",
	Long:  PUSH_COMMAND_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("push called")
	},
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull assets from Github",
	Long:  PULL_COMMAND_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pull called")
	},
}

var drawCmd = &cobra.Command{
	Use:   "draw",
	Short: "Uses go-figure to generate ASCII text",
	Long:  DRAW_COMMAND_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("draw called")
		figure.Draw("anvil", args[0])
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(drawCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
