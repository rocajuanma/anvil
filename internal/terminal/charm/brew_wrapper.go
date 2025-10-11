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
)

// BrewSpinner provides a convenient wrapper for brew operations with spinners
type BrewSpinner struct {
	spinner *Spinner
}

// NewBrewSpinner creates a new brew operation spinner
func NewBrewSpinner() *BrewSpinner {
	return &BrewSpinner{}
}

// InstallPackage shows a spinner while installing a package
func (bs *BrewSpinner) InstallPackage(packageName string, installFunc func() error) error {
	bs.spinner = NewDotsSpinner(fmt.Sprintf("Installing %s", packageName))
	bs.spinner.Start()

	err := installFunc()

	if err != nil {
		bs.spinner.Error(fmt.Sprintf("Failed to install %s", packageName))
		return err
	}

	bs.spinner.Success(fmt.Sprintf("%s installed successfully", packageName))
	return nil
}

// UpdateBrew shows a spinner while updating brew
func (bs *BrewSpinner) UpdateBrew(updateFunc func() error) error {
	bs.spinner = NewDotsSpinner("Updating Homebrew")
	bs.spinner.Start()

	err := updateFunc()

	if err != nil {
		bs.spinner.Error("Failed to update Homebrew")
		return err
	}

	bs.spinner.Success("Homebrew updated successfully")
	return nil
}

// SearchPackage shows a spinner while searching for a package
func (bs *BrewSpinner) SearchPackage(packageName string, searchFunc func() error) error {
	bs.spinner = NewCircleSpinner(fmt.Sprintf("Searching for %s", packageName))
	bs.spinner.Start()

	err := searchFunc()

	if err != nil {
		bs.spinner.Error(fmt.Sprintf("Failed to search for %s", packageName))
		return err
	}

	bs.spinner.Success(fmt.Sprintf("Found %s", packageName))
	return nil
}

// CheckAvailability shows a spinner while checking package availability
func (bs *BrewSpinner) CheckAvailability(packageName string, checkFunc func() (bool, error)) (bool, error) {
	bs.spinner = NewLineSpinner(fmt.Sprintf("Checking %s", packageName))
	bs.spinner.Start()

	available, err := checkFunc()

	if err != nil {
		bs.spinner.Error(fmt.Sprintf("Failed to check %s", packageName))
		return false, err
	}

	if available {
		bs.spinner.Success(fmt.Sprintf("%s is available", packageName))
	} else {
		bs.spinner.Warning(fmt.Sprintf("%s is not installed", packageName))
	}

	return available, nil
}

// InstallHomebrew shows a spinner while installing Homebrew itself
func (bs *BrewSpinner) InstallHomebrew(installFunc func() error) error {
	bs.spinner = NewDotsSpinner("Installing Homebrew (this may take a few minutes)")
	bs.spinner.Start()

	err := installFunc()

	if err != nil {
		bs.spinner.Error("Failed to install Homebrew")
		return err
	}

	bs.spinner.Success("Homebrew installed successfully")
	return nil
}
