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

type RepositoryInterface interface {
	//	List()
	Get(name string) (result *models.RepoRecord, err error)
	//Put()
	//Delete()
}

type repository struct {
	client  rest.Interface
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
		//VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}
