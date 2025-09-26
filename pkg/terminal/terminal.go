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

package terminal

// Color constants for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

var (
	// outputColors is a map of output levels to their corresponding colors
	outputColors = map[OutputLevel]string{
		LevelHeader:  ColorCyan,
		LevelStage:   ColorBlue,
		LevelSuccess: ColorGreen,
		LevelError:   ColorRed,
		LevelWarning: ColorYellow,
		LevelInfo:    "",
	}

	// outputEmojis is a map of output levels to their corresponding emojis
	outputEmojis = map[OutputLevel]string{
		LevelHeader:  "",
		LevelStage:   "🔧 ",
		LevelSuccess: "✅ ",
		LevelError:   "❌ ",
		LevelWarning: "⚠️  ",
		LevelInfo:    "",
	}

	// outputPrefixes is a map of output levels to their corresponding prefixes
	outputPrefixes = map[OutputLevel]string{
		LevelHeader:  headerFormat,
		LevelStage:   "[STAGE] ",
		LevelSuccess: "[SUCCESS] ",
		LevelError:   "[ERROR] ",
		LevelWarning: "[WARNING] ",
		LevelInfo:    "",
	}

	coloredHeaderFormat = "\n%s%s=== %s ===%s\n"
	headerFormat        = "\n=== %s ===\n"
)
