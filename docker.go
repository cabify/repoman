package repoman

import (
	"path"

	"github.com/magefile/mage/sh"
)

func dockerBuildProjectRepos() error {
	for _, repo := range config.Projects {
		if err := dockerBuildRepo(".", repo); err != nil {
			return err
		}
	}
	return nil
}

func dockerBuildGroupRepos() error {
	for group, repos := range config.Groups {
		for _, repo := range repos {
			if err := dockerBuildRepo(group, repo); err != nil {
				return err
			}
		}
	}
	return nil
}

func dockerBuildRepo(group string, r *Repo) error {
	if r.Docker == "" {
		return nil
	}
	p := path.Join(group, r.Name)
	return sh.Run("docker", "build", "-t", r.Docker, p)
}
