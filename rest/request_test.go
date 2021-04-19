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

package rest

import (
	"net/url"
	"testing"
)

func TestNewRequestSetsAccept(t *testing.T) {
	r := NewRequest(nil, "get", &url.URL{Path: "/path/"}, "", ContentConfig{}, nil, 0)
	if r.headers.Get("Accept") != "" {
		t.Errorf("unexpected headers: %#v", r.headers)
	}
	r = NewRequest(nil, "get", &url.URL{Path: "/path/"}, "", ContentConfig{ContentType: "application/other"}, nil, 0)
	if r.headers.Get("Accept") != "application/other, */*" {
		t.Errorf("unexpected headers: %#v", r.headers)
	}
}

func TestRequestAbsPathPreservesTrailingSlash(t *testing.T) {
	r := (&Request{baseURL: &url.URL{}}).AbsPath("/foo/")
	if s := r.URL().String(); s != "/foo/" {
		t.Errorf("trailing slash should be preserved: %s", s)
	}

	r = (&Request{baseURL: &url.URL{}}).AbsPath("/foo/")
	if s := r.URL().String(); s != "/foo/" {
		t.Errorf("trailing slash should be preserved: %s", s)
	}
}

func TestRequestAbsPathJoins(t *testing.T) {
	r := (&Request{baseURL: &url.URL{}}).AbsPath("foo/bar", "baz")
	if s := r.URL().String(); s != "foo/bar/baz" {
		t.Errorf("trailing slash should be preserved: %s", s)
	}
}

func TestRequestSetsNamespace(t *testing.T) {
	r := (&Request{
		baseURL: &url.URL{
			Path: "/",
		},
	}).Project("foo")
	if r.project == "" {
		t.Errorf("namespace should be set: %#v", r)
	}

	if s := r.URL().String(); s != "projects/foo" {
		t.Errorf("namespace should be in path: %s", s)
	}
}

func TestRequestOrdersNamespaceInPath(t *testing.T) {
	r := (&Request{
		baseURL:    &url.URL{},
		pathPrefix: "/test/",
	}).Name("bar").Resource("baz").Project("foo")
	if s := r.URL().String(); s != "/test/projects/foo/baz/bar" {
		t.Errorf("namespace should be in order in path: %s", s)
	}
}
