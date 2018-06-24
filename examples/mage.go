// +build mage

package main

import (
	"github.com/cabify/repoman"
	"github.com/magefile/mage/mg"
)

// Clone and Pull all the project and group repositories.
func Pull() {
	mg.Deps(repoman.Pull)
}

// Run `go get` in each Go repository.
func Get() {
	mg.Deps(repoman.Get)
}

// Determine the Git status of each repository.
func Status() {
	mg.Deps(repoman.Status)
}
