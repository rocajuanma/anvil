package tools

import (
	"fmt"
	"sort"

	"github.com/rocajuanma/anvil/internal/config"
	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/anvil/internal/errors"
)

// LoadAndPrepareAppData loads all application data and prepares it for rendering
// This function is copied from the install package to maintain consistency
func LoadAndPrepareAppData() (groups map[string][]string, builtInGroupNames []string, customGroupNames []string, installedApps []string, err error) {
	// Load groups from config
	groups, err = config.GetAvailableGroups()
	if err != nil {
		err = errors.NewConfigurationError(constants.OpShow, "load-data",
			fmt.Errorf("failed to load groups: %w", err))
		return
	}

	// Get built-in group names
	builtInGroupNames = config.GetBuiltInGroups()

	// Extract and sort custom group names
	for groupName := range groups {
		if !config.IsBuiltInGroup(groupName) {
			customGroupNames = append(customGroupNames, groupName)
		}
	}
	sort.Strings(customGroupNames)

	// Load and sort installed apps
	installedApps, err = config.GetInstalledApps()
	if err != nil {
		// Don't fail on installed apps error, just log warning
		getOutputHandler().PrintWarning("Failed to load installed apps: %v", err)
		installedApps = []string{}
		err = nil // Reset error since we can continue
	} else {
		sort.Strings(installedApps)
	}

	return
}
