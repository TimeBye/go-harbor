package harbor

import (
	"fmt"
	"github.com/goharbor/harbor/src/common/models"
	"github.com/parnurzeal/gorequest"
	"time"
)

// VulnerabilityItem is an item in the vulnerability result returned by vulnerability details API.
type VulnerabilityItem struct {
	ID          string `json:"id"`
	Severity    int64  `json:"severity"`
	Pkg         string `json:"package"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Fixed       string `json:"fixedVersion,omitempty"`
}

type RepoResp struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	ProjectID    int64     `json:"project_id"`
	Description  string    `json:"description"`
	PullCount    int64     `json:"pull_count"`
	StarCount    int64     `json:"star_count"`
	TagsCount    int64     `json:"tags_count"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

// RepoRecord holds the record of an repository in DB, all the infors are from the registry notification event.
type RepoRecord struct {
	*models.RepoRecord
}

type cfg struct {
	Labels map[string]string `json:"labels"`
}

//ComponentsOverview has the total number and a list of components number of different serverity level.
type ComponentsOverview struct {
	Total   int                        `json:"total"`
	Summary []*ComponentsOverviewEntry `json:"summary"`
}

//ComponentsOverviewEntry ...
type ComponentsOverviewEntry struct {
	Sev   int `json:"severity"`
	Count int `json:"count"`
}

//ImgScanOverview mapped to a record of image scan overview.
type ImgScanOverview struct {
	ID              int64               `json:"-"`
	Digest          string              `json:"image_digest"`
	Status          string              `json:"scan_status"`
	JobID           int64               `json:"job_id"`
	Sev             int                 `json:"severity"`
	CompOverviewStr string              `json:"-"`
	CompOverview    *ComponentsOverview `json:"components,omitempty"`
	DetailsKey      string              `json:"details_key"`
	CreationTime    time.Time           `json:"creation_time,omitempty"`
	UpdateTime      time.Time           `json:"update_time,omitempty"`
}

type tagDetail struct {
	Digest        string    `json:"digest"`
	Name          string    `json:"name"`
	Size          int64     `json:"size"`
	Architecture  string    `json:"architecture"`
	OS            string    `json:"os"`
	DockerVersion string    `json:"docker_version"`
	Author        string    `json:"author"`
	Created       time.Time `json:"created"`
	Config        *cfg      `json:"config"`
}

type Signature struct {
	Tag    string            `json:"tag"`
	Hashes map[string][]byte `json:"hashes"`
}

type TagResp struct {
	tagDetail
	Signature    *Signature       `json:"signature"`
	ScanOverview *ImgScanOverview `json:"scan_overview,omitempty"`
}

// RepositoriesService handles communication with the user related methods of
// the Harbor API.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L891
type RepositoriesService struct {
	client *Client
}

type ListRepositoriesOption struct {
	ListOptions
	ProjectId   int64  `url:"project_id,omitempty" json:"project_id,omitempty"`
	ProjectName string `url:"project_name,omitempty" json:"project_name,omitempty"`
	Q           string `url:"q,omitempty" json:"q,omitempty"`
	Sort        string `url:"sort,omitempty" json:"sort,omitempty"`
}

type ManifestResp struct {
	Manifest interface{} `json:"manifest"`
	Config   interface{} `json:"config,omitempty" `
}

// ListRepository Get repositories accompany with relevant project and repo name.
//
// This endpoint let user search repositories accompanying
// with relevant project ID and repo name.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L892
func (s *RepositoriesService) ListRepository(opt *ListRepositoriesOption) ([]RepoRecord, *gorequest.Response, []error) {
	var v []RepoRecord
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf("/projects/%s/repositories", opt.ProjectName)).
		Query(*opt).
		EndStruct(&v)
	return v, &resp, errs
}

// Delete a repository.
//
// This endpoint let user delete a repository with name.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L948
func (s *RepositoriesService) DeleteRepository(repoName string) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.DELETE, fmt.Sprintf("repositories/%s", repoName)).
		End()
	return &resp, errs
}

type RepositoryDescription struct {
	Description string `url:"description,omitempty" json:"description,omitempty"`
}

// Update description of the repository.
//
// This endpoint is used to update description of the repository.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L971
func (s *RepositoriesService) UpdateRepository(repoName string, d RepositoryDescription) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.PUT, fmt.Sprintf("repositories/%s", repoName)).
		Send(d).
		End()
	return &resp, errs
}

// Get the tag of the repository.
//
// This endpoint aims to retrieve the tag of the repository.
// If deployed with Notary, the signature property of
// response represents whether the image is singed or not.
// If the property is null, the image is unsigned.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L999
func (s *RepositoriesService) GetRepositoryTag(repoName, tag string) (TagResp, *gorequest.Response, []error) {
	var v TagResp
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf("repositories/%s/tags/%s", repoName, tag)).
		EndStruct(&v)
	return v, &resp, errs
}

// Delete a tag in a repository.
//
// This endpoint let user delete tags with repo name and tag.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1025
func (s *RepositoriesService) DeleteRepositoryTag(repoName, tag string) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.DELETE, fmt.Sprintf("repositories/%s/tags/%s", repoName, tag)).
		End()
	return &resp, errs
}

// Get tags of a relevant repository.
//
// This endpoint aims to retrieve tags from a relevant repository.
// If deployed with Notary, the signature property of
// response represents whether the image is singed or not.
// If the property is null, the image is unsigned.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1054
func (s *RepositoriesService) ListRepositoryTags(repoName string) ([]TagResp, *gorequest.Response, []error) {
	var v []TagResp
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf("repositories/%s/tags", repoName)).
		EndStruct(&v)
	return v, &resp, errs
}

// Get manifests of a relevant repository.
//
// This endpoint aims to retreive manifests from a relevant repository.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1079
func (s *RepositoriesService) GetRepositoryTagManifests(repoName, tag string, version string) (ManifestResp, *gorequest.Response, []error) {
	var v ManifestResp
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, func() string {
			if version == "" {
				return fmt.Sprintf("repositories/%s/tags/%s/manifest", repoName, tag)
			}
			return fmt.Sprintf("repositories/%s/tags/%s/manifest?version=%s", repoName, tag, version)
		}()).
		EndStruct(&v)
	return v, &resp, errs
}

// Scan the image.
//
// Trigger jobservice to call Clair API to scan the image
// identified by the repo_name and tag.
// Only project admins have permission to scan images under the project.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1113
func (s *RepositoriesService) ScanImage(repoName, tag string) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.POST, fmt.Sprintf("repositories/%s/tags/%s/scan", repoName, tag)).
		End()
	return &resp, errs
}

// Get vulnerability details of the image.
//
// Call Clair API to get the vulnerability based on the previous successful scan.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1177
func (s *RepositoriesService) GetImageDetails(repoName, tag string) ([]VulnerabilityItem, *gorequest.Response, []error) {
	var v []VulnerabilityItem
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf("repositories/%s/tags/%s/vulnerability/details", repoName, tag)).
		EndStruct(&v)
	return v, &resp, errs
}

// Get signature information of a repository.
//
// This endpoint aims to retrieve signature information of a repository, the data is
// from the nested notary instance of Harbor.
// If the repository does not have any signature information in notary, this API will
// return an empty list with response code 200, instead of 404
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1211
func (s *RepositoriesService) GetRepositorySignature(repoName string) ([]Signature, *gorequest.Response, []error) {
	var v []Signature
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf("repositories/%s/signatures", repoName)).
		EndStruct(&v)
	return v, &resp, errs
}

// Get public repositories which are accessed most.
//
// This endpoint aims to let users see the most popular public repositories.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L1241
func (s *RepositoriesService) GetRepositoryTop(top interface{}) ([]RepoResp, *gorequest.Response, []error) {
	var v []RepoResp
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, func() string {
			if t, ok := top.(int); ok {
				return fmt.Sprintf("repositories/top?count=%d", t)
			}
			return fmt.Sprintf("repositories/top")
		}()).
		EndStruct(&v)
	return v, &resp, errs
}
