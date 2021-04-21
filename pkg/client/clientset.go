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

package client

import (
	"fmt"
	project2 "github.com/TimeBye/go-harbor/pkg/project"
	rest2 "github.com/TimeBye/go-harbor/pkg/rest"
	flowcontrol2 "github.com/TimeBye/go-harbor/pkg/rest/util/flowcontrol"
	"github.com/TimeBye/go-harbor/pkg/user"
)

type Interface interface {
	Project() project2.ProjectsInterface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	V2   *project2.ProjectsV2Client
	User *user.UsersClient
}

func NewForConfig(c *rest2.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol2.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	cs := &Clientset{}
	var err error
	cs.V2, err = project2.NewProjectsV1Client(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.User, err = user.NewUsersClient(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return cs, nil
}
