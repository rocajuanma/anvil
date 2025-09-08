/*
Copyright © 2022 Juanma Roca juanmaxroca@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/rocajuanma/anvil/cmd/clean"
	"github.com/rocajuanma/anvil/cmd/config"
	"github.com/rocajuanma/anvil/cmd/doctor"
	"github.com/rocajuanma/anvil/cmd/initcmd"
	"github.com/rocajuanma/anvil/cmd/install"
	"github.com/rocajuanma/anvil/cmd/update"
	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/spf13/cobra"
)

// appVersion holds the version set at build time
var appVersion = "dev"

// SetVersion sets the application version
func SetVersion(v string) {
	appVersion = v
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "anvil",
	Short: "🔥 One CLI to rule them all.",
	Long:  fmt.Sprintf("%s\n\n%s", constants.AnvilLogo, constants.ANVIL_LONG_DESCRIPTION),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if version flag was used
		if version, _ := cmd.Flags().GetBool("version"); version {
			fmt.Printf("Anvil %s\n", appVersion)
			return
		}

		fmt.Println(constants.AnvilLogo)
		fmt.Println()
		fmt.Println("🔥 One CLI to rule them all.")
		fmt.Println()
		fmt.Println("Use 'anvil --help' for usage information.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initcmd.InitCmd)
	rootCmd.AddCommand(install.InstallCmd)
	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(doctor.DoctorCmd)
	rootCmd.AddCommand(clean.CleanCmd)
	rootCmd.AddCommand(update.UpdateCmd)

	// Add version flag
	rootCmd.Flags().BoolP("version", "v", false, "Show version information")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
