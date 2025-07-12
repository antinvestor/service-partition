package models

import (
	"github.com/pitabwire/frame"
)

type Tenant struct {
	frame.BaseModel
	Name        string `gorm:"type:varchar(100);"`
	Description string `gorm:"type:text;"`
	Properties  frame.JSONMap
}

type Partition struct {
	frame.BaseModel
	Name         string `gorm:"type:varchar(100);"`
	Description  string `gorm:"type:text;"`
	ParentID     string `gorm:"type:varchar(50);"`
	ClientSecret string `gorm:"type:varchar(250);"`
	Properties   frame.JSONMap
	State        int32
}

type PartitionRole struct {
	frame.BaseModel
	Name       string `gorm:"type:varchar(100);"`
	Properties frame.JSONMap
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
