package models

import (
	"github.com/pitabwire/frame"
	"gorm.io/datatypes"
)

type Tenant struct {
	frame.BaseModel
	Name        string `gorm:"type:varchar(100);"`
	Description string `gorm:"type:text;"`
	Properties  datatypes.JSONMap
}

type Partition struct {
	frame.BaseModel
	Name        string `gorm:"type:varchar(100);"`
	Description string `gorm:"type:text;"`
	ParentID    string `gorm:"type:varchar(50);"`
	ClientID    string `gorm:"type:varchar(50);"`
	Properties  datatypes.JSONMap
	State       int32
}

type PartitionRole struct {
	frame.BaseModel
	Name       string `gorm:"type:varchar(100);"`
	Properties datatypes.JSONMap
}

type Page struct {
	frame.BaseModel
	Name  string `gorm:"type:varchar(50);"`
	HTML  string `gorm:"type:text;"`
	State int32
}

type Access struct {
	frame.BaseModel
	ProfileID string `gorm:"type:varchar(50);"`
	State     int32
}

type AccessRole struct {
	frame.BaseModel
	AccessID        string `gorm:"type:varchar(50);"`
	PartitionRoleID string `gorm:"type:varchar(50);"`
}
