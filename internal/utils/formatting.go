package utils

import (
	"fmt"

	"github.com/0xjuanma/palantir"
)

func BoldText(text string, color string) string {
	if color == "" {
		color = palantir.ColorBold
	}
	return fmt.Sprintf("%s%s%s", color, text, palantir.ColorReset)
}

func ColorSectionHeader(text string) string {
	return fmt.Sprintf("%s%s%s", palantir.ColorBold+palantir.ColorCyan, text, palantir.ColorReset)
}

func ColorAppName(text string) string {
	return fmt.Sprintf("%s%s%s", palantir.ColorGreen, text, palantir.ColorReset)
}

func ColorGroupNameWithIcon(text string) string {
	return fmt.Sprintf("%s 📁", BoldText(text, ""))
}

func ColoredName(text string, color string) string {
	return BoldText(text, color)
}
