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
	"github.com/TimeBye/go-harbor/pkg/project/options"
	rest2 "github.com/TimeBye/go-harbor/pkg/rest"
	"github.com/goharbor/harbor/src/pkg/project/models"
)

// ProjectsInterface holds the methods that discover server-supported API groups,
// versions and resources.
type ProjectsInterface interface {
	RepositoryInterface
}

// ProjectsV2Client is used to interact with features provided by the admissionregistration.k8s.io group.
type ProjectsV2Client struct {
	restClient rest2.Interface
}

func (p *ProjectsV2Client) Get(name string) (result *models.Project, err error) {
	result = &models.Project{}
	err = p.restClient.Get().
		Resource("projects").
		Name(name).
		Do().
		Into(result)
	return
}

func (p *ProjectsV2Client) List(query *options.ProjectsListOptions) (results *[]models.Project, err error) {
	results = &[]models.Project{}
	err = p.restClient.List().
		Resource("projects").
		Params(*query).
		Do().
		Into(results)
	return
}

func (p *ProjectsV2Client) Delete(name string) (err error) {
	err = p.restClient.Delete().
		Resource("projects").
		Name(name).
		Do().
		Error()
	return
}

func NewProjectsV1Client(restClient *rest2.Config) (*ProjectsV2Client, error) {
	client, err := rest2.RESTClientFor(restClient)
	if err != nil {
		return nil, err
	}
	return &ProjectsV2Client{restClient: client}, nil
}

func (p *ProjectsV2Client) Repositories(project string) RepositoryInterface {
	return newRepositories(p, project)
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (p *ProjectsV2Client) RESTClient() rest2.Interface {
	if p == nil {
		return nil
	}
	return p.restClient
}
