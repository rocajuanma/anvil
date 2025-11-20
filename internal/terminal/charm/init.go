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
	"github.com/0xjuanma/palantir"
)

var (
	// globalCharmHandler is the global enhanced output handler
	globalCharmHandler palantir.OutputHandler
)

// InitCharmOutput initializes the enhanced Charm output handler globally
// Call this once at the start of your application
func InitCharmOutput() {
	globalCharmHandler = NewCharmOutputHandler()
	palantir.SetGlobalOutputHandler(globalCharmHandler)
}

// GetCharmHandler returns the global Charm output handler
func GetCharmHandler() palantir.OutputHandler {
	if globalCharmHandler == nil {
		InitCharmOutput()
	}
	return globalCharmHandler
}

// IsCharmEnabled checks if Charm output is currently enabled
func IsCharmEnabled() bool {
	return globalCharmHandler != nil
}

