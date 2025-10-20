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

package utils

import (
	"fmt"
	"strings"

	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/palantir"
)

// AppTreeNode represents a node in the applications tree
type AppTreeNode struct {
	Name     string
	IsGroup  bool
	Apps     []string
	Children []*AppTreeNode
}

// RenderListView renders applications in a flat list format
func RenderListView(groups map[string][]string, builtInGroupNames []string, customGroupNames []string, installedApps []string) string {
	var content strings.Builder
	content.WriteString("\n")

	// Show built-in groups first
	content.WriteString(ColorSectionHeader("Built-in Groups") + "\n\n")
	for _, groupName := range builtInGroupNames {
		if tools, exists := groups[groupName]; exists {
			content.WriteString(fmt.Sprintf("  %s  %s\n", ColorGroupNameWithIcon(groupName), strings.Join(tools, ", ")))
		}
	}

	// Show custom groups
	if len(customGroupNames) > 0 {
		content.WriteString("\n" + ColorSectionHeader("Custom Groups") + "\n\n")
		for _, groupName := range customGroupNames {
			content.WriteString(fmt.Sprintf("  %s  %s\n", ColorGroupNameWithIcon(groupName), strings.Join(groups[groupName], ", ")))
		}
	} else {
		content.WriteString(fmt.Sprintf("\n%sNo custom groups defined%s\n", palantir.ColorBold+palantir.ColorYellow, palantir.ColorReset))
		content.WriteString(fmt.Sprintf("  Add custom groups in ~/%s/%s\n", constants.AnvilConfigDir, constants.ConfigFileName))
	}

	// Show individually tracked installed apps
	if len(installedApps) > 0 {
		content.WriteString("\n" + ColorSectionHeader("Individually Tracked Apps") + "\n\n")
		for _, app := range installedApps {
			content.WriteString(fmt.Sprintf("  %s\n", ColorAppName(app)))
		}
	}

	content.WriteString("\n")
	return content.String()
}

// RenderTreeView renders applications in a hierarchical tree format
func RenderTreeView(groups map[string][]string, builtInGroupNames []string, customGroupNames []string, installedApps []string) string {
	// Create root node
	root := &AppTreeNode{
		Name:     "Applications",
		IsGroup:  false,
		Children: []*AppTreeNode{},
	}

	// Add built-in groups section
	if len(builtInGroupNames) > 0 {
		builtInNode := &AppTreeNode{
			Name:     "Built-in Groups",
			IsGroup:  false,
			Children: []*AppTreeNode{},
		}

		for _, groupName := range builtInGroupNames {
			if tools, exists := groups[groupName]; exists {
				groupNode := &AppTreeNode{
					Name:    groupName,
					IsGroup: true,
					Apps:    tools,
				}
				builtInNode.Children = append(builtInNode.Children, groupNode)
			}
		}

		if len(builtInNode.Children) > 0 {
			root.Children = append(root.Children, builtInNode)
		}
	}

	// Add custom groups section
	if len(customGroupNames) > 0 {
		customNode := &AppTreeNode{
			Name:     "Custom Groups",
			IsGroup:  false,
			Children: []*AppTreeNode{},
		}

		for _, groupName := range customGroupNames {
			groupNode := &AppTreeNode{
				Name:    groupName,
				IsGroup: true,
				Apps:    groups[groupName],
			}
			customNode.Children = append(customNode.Children, groupNode)
		}

		root.Children = append(root.Children, customNode)
	}

	// Add individually tracked apps section
	if len(installedApps) > 0 {
		individualNode := &AppTreeNode{
			Name:     "Individually Tracked Apps",
			IsGroup:  false,
			Children: []*AppTreeNode{},
		}

		for _, appName := range installedApps {
			appNode := &AppTreeNode{
				Name:    appName,
				IsGroup: false,
			}
			individualNode.Children = append(individualNode.Children, appNode)
		}

		root.Children = append(root.Children, individualNode)
	}

	// Build tree content
	var content strings.Builder
	content.WriteString("\n")
	buildTreeString(&content, root, "", true, true)
	content.WriteString("\n")

	return content.String()
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
