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

	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/palantir"
)

// renderListView renders applications in a flat list format
func renderListView(groups map[string][]string, builtInGroupNames []string, customGroupNames []string, installedApps []string) string {
	var content strings.Builder
	content.WriteString("\n")

	// Show built-in groups first
	content.WriteString(colorSectionHeader("Built-in Groups") + "\n\n")
	for _, groupName := range builtInGroupNames {
		if tools, exists := groups[groupName]; exists {
			content.WriteString(fmt.Sprintf("  %s  %s\n", colorGroupNameWithIcon(groupName), strings.Join(tools, ", ")))
		}
	}

	// Show custom groups
	if len(customGroupNames) > 0 {
		content.WriteString("\n" + colorSectionHeader("Custom Groups") + "\n\n")
		for _, groupName := range customGroupNames {
			content.WriteString(fmt.Sprintf("  %s  %s\n", colorGroupNameWithIcon(groupName), strings.Join(groups[groupName], ", ")))
		}
	} else {
		content.WriteString(fmt.Sprintf("\n%sNo custom groups defined%s\n", palantir.ColorBold+palantir.ColorYellow, palantir.ColorReset))
		content.WriteString(fmt.Sprintf("  Add custom groups in ~/%s/%s\n", constants.AnvilConfigDir, constants.ConfigFileName))
	}

	// Show individually tracked installed apps
	if len(installedApps) > 0 {
		content.WriteString("\n" + colorSectionHeader("Individually Tracked Apps") + "\n\n")
		for _, app := range installedApps {
			content.WriteString(fmt.Sprintf("  %s\n", colorAppName(app)))
		}
	}

	content.WriteString("\n")
	return content.String()
}

// renderTreeView renders applications in a hierarchical tree format
func renderTreeView(groups map[string][]string, builtInGroupNames []string, customGroupNames []string, installedApps []string) string {
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

// Color helper functions for consistent formatting
func colorSectionHeader(text string) string {
	return fmt.Sprintf("%s%s%s", palantir.ColorBold+palantir.ColorCyan, text, palantir.ColorReset)
}

func colorBoldText(text string) string {
	return fmt.Sprintf("%s%s%s", palantir.ColorBold, text, palantir.ColorReset)
}

func colorAppName(text string) string {
	return fmt.Sprintf("%s%s%s", palantir.ColorGreen, text, palantir.ColorReset)
}

func colorGroupNameWithIcon(text string) string {
	return fmt.Sprintf("%s 📁", colorBoldText(text))
}
