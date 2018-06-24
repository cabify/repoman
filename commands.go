package repoman

import (
	"path"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Clone and Pull all the repositories.
func Pull() error {
	mg.Deps(parseConfig)

	// Projects
	mg.Deps(cloneProjectRepos)
	mg.Deps(pullProjectRepos)

	// Groups
	mg.Deps(mkGroupDirs)
	mg.Deps(cloneGroupRepos)
	mg.Deps(pullGroupRepos)

	return nil
}

// Run `go get` in each Go repository.
func Get() error {
	mg.Deps(parseConfig)

	mg.Deps(goGetProjectRepos)

	mg.Deps(goGetGroupRepos)

	return nil
}

// Provide the status of each sub git repository.
func Status() error {
	mg.Deps(parseConfig)
	return sh.RunV("/bin/sh", path.Join(config.Gopath, "src/github.com/cabify/repoman/scripts/mgitstatus.sh"))
}
