package repoman

import (
	"fmt"
	"path"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Pull clones and pulls all the repositories.
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

// Get runs `go get` in each Go repository.
func Get() error {
	mg.Deps(parseConfig)

	mg.Deps(goGetProjectRepos)

	mg.Deps(goGetGroupRepos)

	return nil
}

// Status provides the status of each sub git repository.
func Status() error {
	mg.Deps(parseConfig)
	cmd := []string{
		path.Join(config.Gopath, "src/github.com/cabify/repoman/scripts/mgitstatus.sh"),
	}
	if config.StatusDepth != 0 {
		cmd = append(cmd,
			".",
			fmt.Sprintf("%d", config.StatusDepth),
		)
	}
	return sh.RunV("/bin/sh", cmd...)
}

// DockerBuild attempts to build docker images of each project service.
func DockerBuild() {
	mg.Deps(parseConfig)
	mg.Deps(dockerBuildProjectRepos)
	mg.Deps(dockerBuildGroupRepos)
}
