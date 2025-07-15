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

package draw

import (
	"fmt"

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/figure"
	"github.com/spf13/cobra"
)

var DrawCmd = &cobra.Command{
	Use:   "draw",
	Short: "Uses go-figure to generate ASCII text",
	Long:  constants.DRAW_COMMAND_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("draw called")
		figure.Draw("anvil", args[0])
	},
}

func init() {}
