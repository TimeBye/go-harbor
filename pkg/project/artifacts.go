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
	"fmt"
	"github.com/TimeBye/go-harbor/pkg/model"
	rest2 "github.com/TimeBye/go-harbor/pkg/rest"
)

type ArtifactInterface interface {
	Get(name string) (result *model.Artifact, err error)
	Delete(name string) (err error)
	List(query *model.Query) (result *[]model.Artifact, err error)
}

type artifact struct {
	client     rest2.Interface
	project    string
	repository string
}

// newArtifacts returns a ConfigMaps
func newArtifacts(rest rest2.Interface, project, repository string) *artifact {
	return &artifact{
		client:     rest,
		project:    project,
		repository: repository,
	}
}

func (r *artifact) Get(name string) (result *model.Artifact, err error) {
	result = &model.Artifact{}
	err = r.client.Get().
		Project(r.project).
		Resource("repositories").
		Name(r.repository).
		Suffix(fmt.Sprintf("/artifacts/%s", name)).
		Do().
		Into(result)
	return
}

func (r *artifact) List(query *model.Query) (result *[]model.Artifact, err error) {
	result = &[]model.Artifact{}
	err = r.client.Get().
		Project(r.project).
		Resource("repositories").
		Name(r.repository).
		Suffix("/artifacts").
		Params(*query).
		Do().
		Into(result)
	return
}

func (r *artifact) Delete(name string) (err error) {
	err = r.client.Delete().
		Project(r.project).
		Resource("repositories").
		Name(r.repository).
		Suffix(fmt.Sprintf("/artifacts/%s", name)).
		Do().
		Error()
	return
}
