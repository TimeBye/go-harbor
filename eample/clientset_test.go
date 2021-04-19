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

package eample

import "testing"

func TestNewClientSet(t *testing.T) {
	clientSet, err := NewClientSet("harbor.cloud2go.cn")
	if err != nil {
		t.Errorf("get client set error:%v", err)
	}
	_, err = clientSet.ProjectV1Client.Repositories("xx").Get("xx")
	if err != nil {
		t.Errorf("get repositories set error:%+v", err)
	}
}
