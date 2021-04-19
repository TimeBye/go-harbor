/*
Copyright 2020 The go-harbor Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
*/

package project

import (
	"github.com/TimeBye/go-harbor/rest"
	"github.com/goharbor/harbor/src/common/models"
)

// ProjectsInterface holds the methods that discover server-supported API groups,
// versions and resources.
type ProjectsInterface interface {
	RepositoryInterface
}

// ProjectsV1Client is used to interact with features provided by the admissionregistration.k8s.io group.
type ProjectsV1Client struct {
	restClient rest.Interface
}

func (p *ProjectsV1Client) Get(name string) (result *models.RepoRecord, err error) {
	panic("implement me")
}

func NewProjectsV1Client(restClient *rest.Config) (*ProjectsV1Client, error) {
	client, err := rest.RESTClientFor(restClient)
	if err != nil {
		return nil, err
	}
	return &ProjectsV1Client{restClient: client}, nil
}

func (p *ProjectsV1Client) Repositories(project string) RepositoryInterface {
	return newRepositories(p, project)
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (p *ProjectsV1Client) RESTClient() rest.Interface {
	if p == nil {
		return nil
	}
	return p.restClient
}
