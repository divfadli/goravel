package models

import (
	"github.com/goravel/framework/database/orm"
)

type FileImage struct {
	IdFileImage uint8  `json:"id_file_image" gorm:"primary_key" column:"id_file_image"`
	Filename    string `json:"filename" gorm:"default:not null" column:"filename"`
	Extension   string `json:"extension" gorm:"default:not null" column:"extension"`
	Url         string `json:"url" gorm:"default:not null" column:"url"`
	orm.Timestamps
}

func (r *FileImage) TableName() string {
	return "public.file_image"
}
