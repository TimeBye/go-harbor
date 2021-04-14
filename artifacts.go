package harbor

import (
	"fmt"
	"github.com/goharbor/harbor/src/pkg/artifact"
	"github.com/parnurzeal/gorequest"
)

type Artifact struct {
	artifact.Artifact
}

type ArtifactOption struct {
	ListOptions
	WithImmutableStatus string `json:"with_immutable_status" url:"with_immutable_status"`
	WithLabel           bool   `json:"with_label" url:"with_label"`
	WithScanOverview    bool   `json:"with_scan_overview" url:"with_scan_overview"`
	WithSignature       bool   `json:"with_signature" url:"with_signature"`
	Q                   string `json:"q" url:"q"`
}

//curl -X GET "http://harbor.cloud2go.cn/api/v2.0/projects/sdc/repositories/csd/artifacts?
//q=csd&page=1&page_size=10&with_tag=true&with_label=true&with_scan_overview=true&with_signature=true&with_immutable_status=true"
//-H "accept: application/json" -H "X-Request-Id: sdc" -H "X-Harbor-CSRF-Token: NttN0AYjQEIvqPw97JlFt0S8h81tFmrJBrP0IiU7omM4hNUUhhnZr5QRVLW4aZGHAhP8sRAmFb2aWubXfRFifQ=="

type ArtifactService struct {
	client *Client
}

//List /projects/{project_name}/repositories/{repository_name}/artifacts
func (as *ArtifactService) List(project string, repositories string, opt *ArtifactOption) ([]Artifact, *gorequest.Response, []error) {
	var artifacts []Artifact
	subPath := fmt.Sprintf("projects/%s/repositories/%s/artifacts", project, repositories)
	resp, _, errs := as.client.
		NewRequest(gorequest.GET, subPath).
		Query(*opt).
		EndStruct(&artifacts)
	return artifacts, &resp, errs
}
