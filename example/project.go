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

package example

import (
	"fmt"
	"github.com/TimeBye/go-harbor"
	"github.com/TimeBye/go-harbor/pkg/model"
	"github.com/TimeBye/go-harbor/pkg/project/options"
)

func Project(host, username, password string) error {
	clientSet, err := harbor.NewClientSet(host, username, password)
	if err != nil {
		return fmt.Errorf("get client set error:%v", err)
	}
	result, err := clientSet.V2.Get("15")
	if err != nil || len(result.Name) == 0 {
		return fmt.Errorf("%v", err)
	}

	query := &options.ProjectsListOptions{
		Query: &model.Query{
			PageSize: 2,
			Q:        "name=aa",
		},
		Name: "aa1",
	}
	result1, err := clientSet.V2.List(query)
	if err != nil || len(*result1) == 0 {
		return fmt.Errorf("%v", err)
	}

	err = clientSet.V2.Delete("11")
	if err != nil || len(*result1) == 0 {
		return fmt.Errorf("%v", err)
	}
	return err
}
