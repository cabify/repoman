# repoman: Repository Manager

Repoman is a package to be used in [Magefiles](https://magefile.org/) so that independent repositories can be grouped together and kept in sync. 

The `config.yml` allows repositories to be specified as:

 * Projects, repositories in the current directory.
 * Groups, subdirectories that group together repositories by a common theme.

Repoman has special support for Go projects. If a project is marked with `go: true` in the configuration, it checkout the repository in your `$GOPATH` and automatically symlink from the configurations directory structure.

Checkout the examples directory in this repository to see a quick example on how to use it.


