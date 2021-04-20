# go-harbor

A Harbor API client enabling Go programs to interact with Harbor in a simple and uniform way

[![GitHub license](https://img.shields.io/github/license/TimeBye/go-harbor.svg)](https://github.com/TimeBye/go-harbor/blob/master/LICENSE)
![travis](https://travis-ci.com/ClareChu/go-harbor.svg?branch=release-2.0.0)

## Coverage

This API client package covers most of the existing Harbor API calls and is updated regularly
to add new and/or missing endpoints. Currently the following services are supported:

- [ ] Users
- [x] Projects
- [x] Repositories
- [x] Artifacts
- [ ] Jobs
- [ ] Policies
- [ ] Targets
- [ ] SystemInfo
- [ ] LDAP
- [ ] Configurations

## Usage

```go
import "github.com/TimeBye/go-harbor"
```

Construct a new Harbor client, then use the various services on the client to
access different parts of the Harbor API. For example, to list all
users:

```go
harborClient := harbor.NewClient(nil, "url","username","password")
statisticMap, _, errs := harborClient.GetStatistics()
```

Some API methods have optional parameters that can be passed. For example,
to list all projects for user "haobor":

```go
harborClient := harbor.NewClient(nil, "url","username","password")
opt := &ListProjectsOptions{Name: "haobor"}
projects, _, err := harborClient.Projects.ListProjects(opt)
```

For complete usage of go-harbor, see the full [package docs](https://godoc.org/github.com/TimeBye/go-harbor).

## ToDo

- The biggest thing this package still needs is tests :disappointed:

## Issues

- If you have an issue: report it on the [issue tracker](https://github.com/TimeBye/go-harbor/issues)
