# go-harbor

A Harbor API client enabling Go programs to interact with Harbor in a simple and uniform way

[![GitHub license](https://img.shields.io/github/license/TimeBye/go-harbor.svg)](https://github.com/TimeBye/go-harbor/blob/master/LICENSE)
![travis](https://travis-ci.com/ClareChu/go-harbor.svg?branch=release-2.0.0)
[![codecov](https://codecov.io/gh/ClareChu/go-collection/branch/master/graph/badge.svg?token=zWaoCNi88E)](https://codecov.io/gh/ClareChu/go-collection)

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

[comment]: <> (Construct a new Harbor client, then use the various services on the client to)

[comment]: <> (access different parts of the Harbor API. For example, to list all)

[comment]: <> (users:)

[comment]: <> (```go)

[comment]: <> (harborClient, err := harbor.NewClientSet&#40;"host", "username", "password"&#41;)

[comment]: <> (if err != nil {)

[comment]: <> (	panic&#40;err&#41;)

[comment]: <> (})

[comment]: <> (query := model.Query{})

[comment]: <> (projects, err := harborClient.V2.List&#40;&query&#41;)

[comment]: <> (```)

Some API methods have optional parameters that can be passed. For example,
to list all projects for user "haobor":

```go
harborClient, err := harbor.NewClientSet("host", "username", "password")
if err != nil {
    panic(err)
}
query := options.ProjectsListOptions{}
projects, err := harborClient.V2.List(&query)
```

For complete usage of go-harbor, see the full [package docs](https://godoc.org/github.com/TimeBye/go-harbor).

## ToDo

- The biggest thing this package still needs is tests :disappointed:

## Issues

- If you have an issue: report it on the [issue tracker](https://github.com/TimeBye/go-harbor/issues)
