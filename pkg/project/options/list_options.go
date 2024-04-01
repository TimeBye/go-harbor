package options

import "github.com/TimeBye/go-harbor/pkg/model"

type ProjectsListOptions struct {
	*model.Query
	// Name The name of project.
	Name string `json:"name,omitempty"`
	// The project is public or private.
	//
	Public bool `json:"public,omitempty"`
	//Owner The name of project owner.
	Owner string `json:"owner,omitempty"`
}

type RepositoriesListOptions struct {
	*model.Query
	// ProjectName The name of the project
	ProjectName string `json:"project_name,omitempty"`
}

type ArtifactsListOptions struct {
	*model.Query
	// ProjectName The name of the project
	ProjectName string `json:"project_name,omitempty"`
	// The name of the repository. If it contains slash, encode it with URL encoding. e.g. a/b -> a%252Fb
	RepositoryName string `json:"repository_name,omitempty"`
	// Specify whether the tags are included inside the returning artifacts
	// Default value : true
	WithTag bool `json:"with_tag,omitempty"`
	//Specify whether the labels are included inside the returning artifacts
	//Default value : false
	WithLabel bool `json:"with_label,omitempty"`
	// Specify whether the scan overview is included inside the returning artifacts
	//Default value : false
	WithScanOverview bool `json:"with_scan_overview,omitempty"`
	// Specify whether the signature is included inside the tags of the returning artifacts. Only works when setting "with_tag=true"
	//Default value : false
	WithSignature bool `json:"with_signature,omitempty"`
	// Specify whether the immutable status is included inside the tags of the returning artifacts. Only works when setting "with_tag=true"
	//Default value : false
	WithImmutableStatus bool `json:"with_immutable_status,omitempty"`
}
