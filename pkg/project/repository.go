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
	"github.com/TimeBye/go-harbor/pkg/model"
	"github.com/TimeBye/go-harbor/pkg/project/options"
	rest2 "github.com/TimeBye/go-harbor/pkg/rest"
)

type RepositoryInterface interface {
	Artifacts(Repository string) *artifact
	List(query *options.RepositoriesListOptions) (result *[]model.Repository, err error)
	Get(name string) (result *model.Repository, err error)
	Delete(name string) (err error)
	//Put()
}

type Repository struct {
	client  rest2.Interface
	project string
}

// newRepositories returns a ConfigMaps
func newRepositories(
	c *ProjectsV2Client,
	project string) *Repository {
	return &Repository{
		client:  c.RESTClient(),
		project: project,
	}
}

func (r *Repository) Get(name string) (result *model.Repository, err error) {
	result = &model.Repository{}
	err = r.client.Get().
		Project(r.project).
		Resource("repositories").
		Name(name).
		Do().
		Into(result)
	return
}

func (r *Repository) List(query *options.RepositoriesListOptions) (result *[]model.Repository, err error) {
	result = &[]model.Repository{}
	err = r.client.Get().
		Project(r.project).
		Resource("repositories").
		Params(*query).
		Do().
		Into(result)
	return
}

func (r *Repository) Delete(name string) (err error) {
	err = r.client.Delete().
		Project(r.project).
		Resource("repositories").
		Name(name).
		Do().
		Error()
	return
}

func (r *Repository) Artifacts(Repository string) *artifact {
	return newArtifacts(r.client, r.project, Repository)
}
