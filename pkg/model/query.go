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

package model

type Query struct {
	// PageSize The size of per page
	// Default value : 10
	PageSize int64 `json:"page_size,omitempty"`
	//Page  The page number
	// Default value : 1
	Page int64 `json:"page,omitempty"`
	// Query string to query resources. Supported query patterns are "exact match(k=v)",
	//"fuzzy match(k=~v)", "range(k=[min~max])", "list with union releationship(k={v1 v2 v3})"
	//and "list with intersetion relationship(k=(v1 v2 v3))". The value of range and list can be string(enclosed by " or '),
	//integer or time(in format "2020-04-09 02:36:00"). All of these query patterns should be put in the query string "q=xxx"
	//and splitted by ",". e.g. q=k1=v1,k2=~v2,k3=[min~max]
	Q string `json:"q,omitempty"`
	// An unique ID for the request
	RequestId string `json:"X-Request-Id,omitempty"`
}
