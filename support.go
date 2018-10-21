package repoman

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/magefile/mage/sh"
)

func mkGroupDirs() error {
	for group := range config.Groups {
		gpath := groupNameToPath(group)
		if _, err := os.Stat(gpath); os.IsNotExist(err) {
			fmt.Printf("Creating group: %v\n", gpath)
			if err = os.MkdirAll(gpath, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}

func groupNameToPath(group string) string {
	p := strings.Split(group, ".")
	return path.Join(p...)
}

func cloneGroupRepos() error {
	for group, repos := range config.Groups {
		for _, repo := range repos {
			if err := cloneRepo(group, repo); err != nil {
				return err
			}
		}
	}
	return nil
}

func cloneProjectRepos() error {
	for _, repo := range config.Projects {
		if err := cloneRepo(".", repo); err != nil {
			return err
		}
	}
	return nil
}

func pullGroupRepos() error {
	for group, repos := range config.Groups {
		for _, repo := range repos {
			if err := pullRepo(group, repo); err != nil {
				return err
			}
		}
	}
	return nil
}

func pullProjectRepos() error {
	for _, repo := range config.Projects {
		if err := pullRepo(".", repo); err != nil {
			return err
		}
	}
	return nil
}

func cloneRepo(group string, r *Repo) error {
	if r.Go {
		// Go Repos require special treatment to ensure they're cloned into the GOPATH
		p := r.GoOrgPath()
		if _, err := os.Stat(p); os.IsNotExist(err) {
			log.Printf("Creating Go Org Path: %v\n", p)
			if err = os.Mkdir(p, os.ModePerm); err != nil {
				return err
			}
		}
		if _, err := os.Stat(r.GoProjectPath()); os.IsNotExist(err) {
			log.Printf("Cloning Go Repo: %v\n", r.GitRepo())
			err = sh.Run("git", "-C", r.GoOrgPath(), "clone", r.GitRepo())
			if err != nil {
				return err
			}
		}

		gpath := groupNameToPath(group)
		p = path.Join(gpath, r.Name)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			log.Printf("Symlinking Go Repo: %v\n", p)
			if err = os.Symlink(r.GoProjectPath(), p); err != nil {
				return err
			}
		}
	} else {
		// Normal repos are much easier!
		gpath := groupNameToPath(group)
		p := path.Join(gpath, r.Name)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			log.Printf("Cloning Repo %v into %v\n", r.GitRepo(), p)
			err := sh.Run("git", "-C", gpath, "clone", r.GitRepo(), r.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func pullRepo(group string, r *Repo) error {
	gpath := groupNameToPath(group)
	p := path.Join(gpath, r.Name)
	if _, err := os.Stat(p); err != nil {
		return err
	}

	log.Printf("Pulling Repo: %v\n", p)
	sh.Run("git", "-C", p, "pull")

	return nil
}

func goGetProjectRepos() error {
	for _, repo := range config.Projects {
		if repo.Go {
			if err := goGetRepo(".", repo); err != nil {
				return err
			}
		}
	}
	return nil
}

func goGetGroupRepos() error {
	for group, repos := range config.Groups {
		gpath := groupNameToPath(group)
		for _, repo := range repos {
			if repo.Go {
				if err := goGetRepo(gpath, repo); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func goGetRepo(group string, r *Repo) error {
	gpath := groupNameToPath(group)
	p := path.Join(gpath, r.Name)
	if _, err := os.Stat(p); err != nil {
		return err
	}

	// Go doesn't support setting the Working Dir from the command line
	os.Chdir(p)
	defer os.Chdir(config.pwd)

	var err error
	if _, err = os.Stat("Gopkg.toml"); err == nil {
		os.Chdir(r.GoProjectPath()) // Need to be in project base path
		log.Printf("Go 'dep ensure' in Repo: %v\n", p)
		err = sh.Run("dep", "ensure", "-v", "-vendor-only")
	} else {
		log.Printf("Go Get in Repo: %v\n", p)
		err = sh.Run("go", "get", "./...")
	}

	return err
}
