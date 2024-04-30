package model

import (
	"time"
)

// Repository holds the record of an repository in DB, all the infors are from the registry notification event.
type Repository struct {
	Id           int64     `orm:"pk;auto;column(id)" json:"id"`
	RepositoryID int64     `orm:"pk;auto;column(repository_id)" json:"repository_id"`
	Name         string    `orm:"column(name)" json:"name"`
	ProjectID    int64     `orm:"column(project_id)"  json:"project_id"`
	Description  string    `orm:"column(description)" json:"description"`
	PullCount    int64     `orm:"column(pull_count)" json:"pull_count"`
	StarCount    int64     `orm:"column(star_count)" json:"star_count"`
	CreationTime time.Time `orm:"column(creation_time);auto_now_add" json:"creation_time"`
	UpdateTime   time.Time `orm:"column(update_time);auto_now" json:"update_time"`
}
