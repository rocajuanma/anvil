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

package install

import (
	"fmt"
	"strings"

	"github.com/rocajuanma/palantir"
)

// AppTreeNode represents a node in the applications tree
type AppTreeNode struct {
	Name     string
	IsGroup  bool
	Apps     []string
	Children []*AppTreeNode
}

// buildTreeString writes an app tree node to a string builder with ASCII art and colors
func buildTreeString(builder *strings.Builder, node *AppTreeNode, prefix string, isLast bool, isRoot bool) {
	if !isRoot {
		var treeChar string
		if isLast {
			treeChar = "└── "
		} else {
			treeChar = "├── "
		}

		// Color the output based on node type
		var coloredName string
		if node.IsGroup {
			// Groups are colored in bold blue
			coloredName = fmt.Sprintf("%s%s%s 📁 %s", palantir.ColorBold, palantir.ColorBlue, node.Name, palantir.ColorReset)
		} else if len(node.Children) > 0 {
			// Category headers (Built-in Groups, Custom Groups, etc.) in bold cyan
			coloredName = fmt.Sprintf("%s%s%s%s", palantir.ColorBold, palantir.ColorCyan, node.Name, palantir.ColorReset)
		} else {
			// Individual apps in green
			coloredName = fmt.Sprintf("%s%s%s", palantir.ColorGreen, node.Name, palantir.ColorReset)
		}
		builder.WriteString(fmt.Sprintf("%s%s%s\n", prefix, treeChar, coloredName))
	}

	// Write apps within a group
	if node.IsGroup && len(node.Apps) > 0 {
		for i, app := range node.Apps {
			isAppLast := i == len(node.Apps)-1

			// Calculate prefix for app
			var appPrefix string
			if isLast {
				appPrefix = prefix + "    "
			} else {
				appPrefix = prefix + "│   "
			}

			var appTreeChar string
			if isAppLast {
				appTreeChar = "└── "
			} else {
				appTreeChar = "├── "
			}

			// Color individual apps in green
			coloredApp := fmt.Sprintf("%s%s%s", palantir.ColorGreen, app, palantir.ColorReset)
			builder.WriteString(fmt.Sprintf("%s%s%s\n", appPrefix, appTreeChar, coloredApp))
		}
	}

	// Write children
	if node.Children != nil {
		for i, child := range node.Children {
			isChildLast := i == len(node.Children)-1

			// Calculate prefix for child
			var childPrefix string
			if isRoot {
				childPrefix = ""
			} else {
				if isLast {
					childPrefix = prefix + "    "
				} else {
					childPrefix = prefix + "│   "
				}
			}

			buildTreeString(builder, child, childPrefix, isChildLast, false)
		}
	}
}
