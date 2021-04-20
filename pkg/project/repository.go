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
	rest2 "github.com/TimeBye/go-harbor/pkg/rest"
	"github.com/goharbor/harbor/src/common/models"
)

type RepositoryInterface interface {
	List(query *model.Query) (result *[]models.RepoRecord, err error)
	Get(name string) (result *models.RepoRecord, err error)
	Delete(name string) (err error)
	//Put()
}

type repository struct {
	client  rest2.Interface
	project string
}

// newConfigMaps returns a ConfigMaps
func newRepositories(
	c *ProjectsV1Client,
	project string) *repository {
	return &repository{
		client:  c.RESTClient(),
		project: project,
	}
}

func (r *repository) Get(name string) (result *models.RepoRecord, err error) {
	result = &models.RepoRecord{}
	err = r.client.Get().
		Project(r.project).
		Resource("repositories").
		Name(name).
		Do().
		Into(result)
	return
}

func (r *repository) List(query *model.Query) (result *[]models.RepoRecord, err error) {
	result = &[]models.RepoRecord{}
	err = r.client.Get().
		Project(r.project).
		Resource("repositories").
		Params(*query).
		Do().
		Into(result)
	return
}

func (r *repository) Delete(name string) (err error) {
	err = r.client.Delete().
		Project(r.project).
		Resource("repositories").
		Name(name).
		Do().
		Error()
	return
}
