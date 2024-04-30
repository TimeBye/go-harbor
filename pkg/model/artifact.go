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

import (
	"github.com/goharbor/harbor/src/pkg/artifact"
	cmodels "github.com/goharbor/harbor/src/pkg/label/model"
	"github.com/goharbor/harbor/src/pkg/tag/model/tag"
)

type Artifact struct {
	artifact.Artifact
	Tags          []*tag.Tag               `json:"tags"`           // the list of tags that attached to the artifact
	AdditionLinks map[string]*AdditionLink `json:"addition_links"` // the resource link for build history(image), values.yaml(chart), dependency(chart), etc
	Labels        []*cmodels.Label         `json:"labels"`
}

// AdditionLink is a link via that the addition can be fetched
type AdditionLink struct {
	HREF     string `json:"href"`
	Absolute bool   `json:"absolute"` // specify the href is an absolute URL or not
}
