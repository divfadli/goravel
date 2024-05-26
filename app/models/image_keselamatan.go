package models

import (
	"github.com/goravel/framework/database/orm"
)

type ImageKeselamatan struct {
	IDImageKeselamatan    uint8               `json:"id_image_keselamatan" gorm:"primary_key" column:"id_image_keselamatan"`
	FileImageID           int                 `json:"file_image_id" gorm:"default:not null" column:"file_image_id"`
	FileImage             FileImage           `gorm:"foreign_key:FileImageID"`
	KejadianKeselamatanID int                 `json:"kejadian_keselamatan_id" gorm:"default:not null" column:"kejadian_keselamatan_id"`
	KejadianKeselamatan   KejadianKeselamatan `gorm:"foreign_key:KejadianKeselamatanID"`
	orm.Timestamps
}

func (r *ImageKeselamatan) TableName() string {
	return "public.image_keselamatan"
}
