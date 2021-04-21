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

package user

import (
	"github.com/TimeBye/go-harbor/pkg/model"
	rest2 "github.com/TimeBye/go-harbor/pkg/rest"
	"github.com/goharbor/harbor/src/common/models"
)

// UsersInterface holds the methods that discover server-supported API groups,
// versions and resources.
type UsersInterface interface {
}

type UsersClient struct {
	restClient rest2.Interface
}

func NewUsersClient(restClient *rest2.Config) (*UsersClient, error) {
	client, err := rest2.RESTClientFor(restClient)
	if err != nil {
		return nil, err
	}
	return &UsersClient{restClient: client}, nil
}

func (u *UsersClient) Get(name string) (result *models.User, err error) {
	result = &models.User{}
	err = u.restClient.Get().
		Resource("users").
		Name(name).
		Do().
		Into(result)
	return
}

func (u *UsersClient) List(query *model.Query) (results *[]models.User, err error) {
	results = &[]models.User{}
	err = u.restClient.List().
		Resource("users").
		Params(*query).
		Do().
		Into(results)
	return
}

func (u *UsersClient) Delete(name string) (err error) {
	return u.restClient.Delete().
		Resource("users").
		Name(name).
		Do().
		Error()
}
