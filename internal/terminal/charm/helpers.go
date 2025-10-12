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

package charm

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderBox creates a beautiful box around content with max width of 120 chars
func RenderBox(title, content string, borderColor string) string {
	if borderColor == "" {
		borderColor = "#FF6B9D"
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1).
		MarginTop(1).
		MarginBottom(1).
		Width(120)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(borderColor))

	header := titleStyle.Render(title)
	return boxStyle.Render(header + "\n\n" + content)
}

// RenderList creates a styled list of items
func RenderList(items []string, bullet string, color string) string {
	if bullet == "" {
		bullet = "•"
	}
	if color == "" {
		color = "#87CEEB"
	}

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		PaddingLeft(2)

	var result strings.Builder
	for _, item := range items {
		result.WriteString(itemStyle.Render(bullet + " " + item + "\n"))
	}
	return result.String()
}

// RenderTable creates a simple styled table
func RenderTable(headers []string, rows [][]string) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6B9D")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("#FF6B9D")).
		Padding(0, 2)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 2)

	// Render headers
	var result strings.Builder
	headerRow := ""
	for _, h := range headers {
		headerRow += headerStyle.Render(h)
	}
	result.WriteString(headerRow + "\n")

	// Render rows
	for _, row := range rows {
		rowStr := ""
		for _, cell := range row {
			rowStr += cellStyle.Render(cell)
		}
		result.WriteString(rowStr + "\n")
	}

	return result.String()
}

// RenderBanner creates a large stylized banner
func RenderBanner(text string) string {
	bannerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6B9D")).
		Background(lipgloss.Color("#2D2D2D")).
		Padding(1, 4).
		MarginTop(1).
		MarginBottom(1).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#FF6B9D")).
		Align(lipgloss.Center)

	return bannerStyle.Render(text)
}

// RenderKeyValue creates a styled key-value pair
func RenderKeyValue(key, value string) string {
	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D9FF")).
		Width(20).
		Align(lipgloss.Right)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB"))

	return keyStyle.Render(key) + " " + valueStyle.Render(value)
}

// RenderSeparator creates a visual separator line
func RenderSeparator(width int, char string, color string) string {
	if char == "" {
		char = "─"
	}
	if color == "" {
		color = "#666666"
	}

	line := strings.Repeat(char, width)
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		MarginTop(1).
		MarginBottom(1)

	return style.Render(line)
}

// RenderHighlight highlights important text
func RenderHighlight(text string, color string) string {
	if color == "" {
		color = "#FFD700"
	}

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(color)).
		Background(lipgloss.Color("#2D2D2D")).
		Padding(0, 1)

	return style.Render(text)
}

// RenderCode renders text as code
func RenderCode(code string) string {
	codeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C792EA")).
		Background(lipgloss.Color("#1E1E1E")).
		Padding(0, 1).
		Italic(true)

	return codeStyle.Render(code)
}

// RenderQuote renders text as a quote
func RenderQuote(quote string, author string) string {
	quoteStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Italic(true).
		PaddingLeft(4).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#00D9FF"))

	authorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		PaddingLeft(6).
		Italic(true)

	result := quoteStyle.Render(quote)
	if author != "" {
		result += "\n" + authorStyle.Render("— "+author)
	}

	return result
}

// RenderBadge creates a small badge/tag
func RenderBadge(text string, color string) string {
	if color == "" {
		color = "#00D9FF"
	}

	badgeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color(color)).
		Padding(0, 1).
		Bold(true)

	return badgeStyle.Render(text)
}

// RenderSteps creates a numbered step list
func RenderSteps(steps []string) string {
	stepNumberStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D9FF")).
		Background(lipgloss.Color("#2D2D2D")).
		Padding(0, 1).
		MarginRight(1)

	stepTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB"))

	var result strings.Builder
	for i, step := range steps {
		number := stepNumberStyle.Render(fmt.Sprintf("%d", i+1))
		text := stepTextStyle.Render(step)
		result.WriteString(number + " " + text + "\n")
	}

	return result.String()
}

// RenderStatus creates a status indicator
func RenderStatus(status string, isPositive bool) string {
	var color string
	var icon string

	if isPositive {
		color = "#00FF87"
		icon = "●"
	} else {
		color = "#FF5F87"
		icon = "●"
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true)

	return statusStyle.Render(icon + " " + status)
}

// RenderPercentage creates a styled percentage display
func RenderPercentage(value float64) string {
	var color string
	if value >= 80 {
		color = "#00FF87"
	} else if value >= 50 {
		color = "#FFD700"
	} else {
		color = "#FF5F87"
	}

	percentStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(color))

	return percentStyle.Render(fmt.Sprintf("%.1f%%", value))
}
