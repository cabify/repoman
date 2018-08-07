package repoman

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Repo struct {
	Name   string `yaml:"name"` // Allows name to be overridden (might be useful?)
	Github string `yaml:"github"`
	Docker string `yaml:"docker"` // Docker image name
	Go     bool   `yaml:"go"`
}

type Group map[string]*Repo

type Config struct {
	Gopath   string           `yaml:"gotpath"`
	Projects map[string]*Repo `yaml:"projects"`
	Groups   map[string]Group `yaml:"groups"`
	pwd      string           `yaml:"-"` // current working, for reference
}

var config Config

func (r *Repo) GitRepo() string {
	return fmt.Sprintf("git@github.com:%v.git", r.Github)
}

func (r *Repo) GitOrg() string {
	return strings.SplitN(r.Github, "/", 2)[0]
}

func (r *Repo) GitProject() string {
	return strings.SplitN(r.Github, "/", 2)[1]
}

func (r *Repo) GitProvider() string {
	// We only support github.com at the moment!
	return "github.com"
}

func (r *Repo) GoOrgPath() string {
	return path.Join(config.Gopath, "src", r.GitProvider(), r.GitOrg())
}

func (r *Repo) GoProjectPath() string {
	return path.Join(r.GoOrgPath(), r.GitProject())
}

func parseConfig() error {
	if len(config.Groups) > 0 || len(config.Projects) > 0 {
		// Already loaded
		return nil
	}
	return parseConfigAt("./config.yml")
}

func parseConfigAt(file string) error {
	config.Gopath = os.Getenv("GOPATH")
	if config.Gopath == "" {
		return errors.New("Doesn't look like your GOPATH is configured!")
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(data, &config); err != nil {
		return err
	}
	// Copy project and group Repo names
	for name, repo := range config.Projects {
		if repo.Name == "" {
			repo.Name = name
		}
	}
	for _, group := range config.Groups {
		for name, repo := range group {
			if repo.Name == "" {
				repo.Name = name
			}
		}
	}

	config.pwd, err = os.Getwd()
	return err
}
