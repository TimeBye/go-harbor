package harbor

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"time"
)

// ProjectMetadata holds the metadata of a project.
type ProjectMetadata struct {
	ID           int64     `json:"id"`
	ProjectID    int64     `json:"project_id"`
	Name         string    `json:"name"`
	Value        string    `json:"value"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
	Deleted      int       `json:"deleted"`
}

// Project holds the details of a project.
type Project struct {
	ProjectID    int64             `json:"project_id"`
	OwnerID      int               `json:"owner_id"`
	Name         string            `json:"name"`
	CreationTime time.Time         `json:"creation_time"`
	UpdateTime   time.Time         `json:"update_time"`
	Deleted      int               `json:"deleted"`
	OwnerName    string            `json:"owner_name"`
	Togglable    bool              `json:"togglable"`
	Role         int               `json:"current_user_role_id"`
	RepoCount    int64             `json:"repo_count"`
	Metadata     map[string]string `json:"metadata"`
}

// AccessLog holds information about logs which are used to record the actions that user take to the resourses.
type AccessLog struct {
	LogID     int       `json:"log_id"`
	Username  string    `json:"username"`
	ProjectID int64     `json:"project_id"`
	RepoName  string    `json:"repo_name"`
	RepoTag   string    `json:"repo_tag"`
	GUID      string    `json:"guid"`
	Operation string    `json:"operation"`
	OpTime    time.Time `json:"op_time"`
}

// ProjectRequest holds informations that need for creating project API
type ProjectRequest struct {
	Name     string            `url:"name,omitempty" json:"project_name"`
	Public   *int              `url:"public,omitempty" json:"public"` //deprecated, reserved for project creation in replication
	Metadata map[string]string `url:"-" json:"metadata"`
}

type ListProjectsOptions struct {
	ListOptions
	Name   string `url:"name,omitempty" json:"name,omitempty"`
	Public bool   `url:"public,omitempty" json:"public,omitempty"`
	Owner  string `url:"owner,omitempty" json:"owner,omitempty"`
}

//ListLogOptions LogQueryParam is used to set query conditions when listing
// access logs.
type ListLogOptions struct {
	ListOptions
	Username   string     `url:"username,omitempty"`        // the operator's username of the log
	Repository string     `url:"repository,omitempty"`      // repository name
	Tag        string     `url:"tag,omitempty"`             // tag name
	Operations []string   `url:"operation,omitempty"`       // operations
	BeginTime  *time.Time `url:"begin_timestamp,omitempty"` // the time after which the operation is done
	EndTime    *time.Time `url:"end_timestamp,omitempty"`   // the time before which the operation is doen
}

type MemberRequest struct {
	UserName string `json:"username"`
	Roles    []int  `json:"roles"`
}

// ProjectsService handles communication with the user related methods of
// the Harbor API.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L45
type ProjectsService struct {
	client *Client
}

//ListProject List projects
//
// This endpoint returns all projects created by Harbor,
// and can be filtered by project name.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L46
func (s *ProjectsService) ListProject(opt *ListProjectsOptions) ([]Project, *gorequest.Response, []error) {
	var projects []Project
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, "projects").
		Query(*opt).
		EndStruct(&projects)
	return projects, &resp, errs
}

//CheckProject Check if the project name user provided already exists.
//
// This endpoint is used to check if the project name user provided already exist.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L100
func (s *ProjectsService) CheckProject(projectName string) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.HEAD, "projects").
		Query(fmt.Sprintf("project_name=%s", projectName)).
		End()
	return &resp, errs
}

// Create a new project.
//
// This endpoint is for user to create a new project.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L122
func (s *ProjectsService) CreateProject(p ProjectRequest) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.POST, "projects").
		Send(p).
		End()
	return &resp, errs
}

// Return specific project detail information.
//
// This endpoint returns specific project information by project ID.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L149
func (s *ProjectsService) GetProjectByID(pid int64) (Project, *gorequest.Response, []error) {
	var project Project
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf("projects/%d", pid)).
		EndStruct(&project)
	return project, &resp, errs
}

// Update properties for a selected project.
//
// This endpoint is aimed to update the properties of a project.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L171
func (s *ProjectsService) UpdateProject(pid int64, p Project) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.PUT, fmt.Sprintf("projects/%d", pid)).
		Send(p).
		End()
	return &resp, errs
}

// Delete project by projectID.
//
// This endpoint is aimed to delete project by project ID.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L203
func (s *ProjectsService) DeleteProject(pid int64) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.DELETE, fmt.Sprintf("projects/%d", pid)).
		End()
	return &resp, errs
}

// Get access logs accompany with a relevant project.
//
// This endpoint let user search access logs filtered by operations and date time ranges.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L230
func (s *ProjectsService) GetProjectLogByID(pid int64, opt ListLogOptions) ([]AccessLog, *gorequest.Response, []error) {
	var accessLog []AccessLog
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf("projects/%d", pid)).
		Query(opt).
		EndStruct(&accessLog)
	return accessLog, &resp, errs
}

// Get project all metadata.
//
// This endpoint returns metadata of the project specified by project ID.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L307
func (s *ProjectsService) GetProjectMetadataById(pid int64) (map[string]string, *gorequest.Response, []error) {
	var metadata map[string]string
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf("projects/%d", pid)).
		EndStruct(&metadata)
	return metadata, &resp, errs
}

// Add metadata for the project.
//
// This endpoint is aimed to add metadata of a project.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L329
func (s *ProjectsService) AddProjectMetadata(pid int64, metadata map[string]string) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.POST, fmt.Sprintf("projects/%d/metadatas", pid)).
		Send(metadata).
		End()
	return &resp, errs
}

// Get project metadata
//
// This endpoint returns specified metadata of a project.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L364
func (s *ProjectsService) GetProjectMetadata(pid int64, specified string) (map[string]string, *gorequest.Response, []error) {
	var metadata map[string]string
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf("projects/%d/metadatas/%s", pid, specified)).
		EndStruct(&metadata)
	return metadata, &resp, errs
}

// Update metadata of a project.
//
// This endpoint is aimed to update the metadata of a project.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L391
func (s *ProjectsService) UpdateProjectMetadata(pid int64, metadataName string) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.PUT, fmt.Sprintf("projects/%d/%s", pid, metadataName)).
		End()
	return &resp, errs
}

// Delete metadata of a project
//
// This endpoint is aimed to delete metadata of a project.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L422
func (s *ProjectsService) DeleteProjectMetadata(pid int64, metadataName string) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.DELETE, fmt.Sprintf("projects/%d/%s", pid, metadataName)).
		End()
	return &resp, errs
}

// User holds the details of a user.
type User struct {
	UserID       int       `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Realname     string    `json:"realname"`
	Comment      string    `json:"comment"`
	Deleted      int       `json:"deleted"`
	Rolename     string    `json:"role_name"`
	Role         int       `json:"role_id"`
	RoleList     []Role    `json:"role_list"`
	HasAdminRole int       `json:"has_username_role"`
	ResetUUID    string    `json:"reset_uuid"`
	Salt         string    `json:"-"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

// Return a project's relevant role members.
//
// This endpoint is for user to search a specified project’s relevant role members.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L452
func (s *ProjectsService) GetProjectMembers(pid int64) ([]User, *gorequest.Response, []error) {
	var users []User
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf("projects/%d/members", pid)).
		EndStruct(&users)
	return users, &resp, errs
}

// Add project role member accompany with relevant project and user.
//
// This endpoint is for user to add project role member accompany with relevant project and user.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L483
func (s *ProjectsService) AddProjectMember(pid int64, member MemberRequest) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.POST, fmt.Sprintf("projects/%d/metadatas", pid)).
		Send(member).
		End()
	return &resp, errs
}

// Role holds the details of a role.
type Role struct {
	RoleID   int    `json:"role_id"`
	RoleCode string `json:"role_code"`
	Name     string `json:"role_name"`
	RoleMask int    `json:"role_mask"`
}

// Return role members accompany with relevant project and user.
//
// This endpoint is for user to get role members accompany with relevant project and user.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L522
func (s *ProjectsService) GetProjectMemberRole(pid, uid int) (Role, *gorequest.Response, []error) {
	var role Role
	resp, _, errs := s.client.
		NewRequest(gorequest.GET, fmt.Sprintf("projects/%d/members/%d", pid, uid)).
		EndStruct(&role)
	return role, &resp, errs
}

// Update project role members accompany with relevant project and user.
//
// This endpoint is for user to update current project role members accompany with relevant project and user.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L559
func (s *ProjectsService) UpdateProjectMemberRole(pid, uid int, role MemberRequest) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.PUT, fmt.Sprintf("projects/%d/members/%d", pid, uid)).
		Send(role).
		End()
	return &resp, errs
}

// Delete project role members accompany with relevant project and user.
//
// This endpoint is aimed to remove project role members already added to the relevant project and user.
//
// Harbor API docs: https://github.com/vmware/harbor/blob/release-1.4.0/docs/swagger.yaml#L597
func (s *ProjectsService) DeleteProjectMember(pid, uid int) (*gorequest.Response, []error) {
	resp, _, errs := s.client.
		NewRequest(gorequest.DELETE, fmt.Sprintf("projects/%d/members/%d", pid, uid)).
		End()
	return &resp, errs
}
