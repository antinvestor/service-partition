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
	Properties  datatypes.JSONMap
	State       int32
}

type PartitionRole struct {
	frame.BaseModel
	Name       string `gorm:"type:varchar(100);"`
	Partition  Partition
	Properties datatypes.JSONMap
}

type Page struct {
	frame.BaseModel
	Name      string `gorm:"type:varchar(50);"`
	Html      string `gorm:"type:text;"`
	Partition Partition
	State     int32
}

type Access struct {
	frame.BaseModel
	ProfileID string `gorm:"type:varchar(50);"`
	Partition Partition
	State     int32
}

type AccessRole struct {
	frame.BaseModel
	Access          Access
	AccessID        string `gorm:"type:varchar(50);"`
	PartitionRole   PartitionRole
	PartitionRoleID string
}
