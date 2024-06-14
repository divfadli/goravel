package models

import (
	"github.com/goravel/framework/database/orm"
)

type ImageKeselamatan struct {
	IDImageKeselamatan    uint8               `json:"id_image_keselamatan" gorm:"primary_key" column:"id_image_keselamatan"`
	FileImageID           uint8               `json:"file_image_id" gorm:"default:0" column:"file_image_id"`
	FileImage             FileImage           `gorm:"foreignKey:FileImageID;references:IdFileImage"`
	KejadianKeselamatanID uint8               `json:"kejadian_keselamatan_id" gorm:"default:0" column:"kejadian_keselamatan_id"`
	KejadianKeselamatan   KejadianKeselamatan `gorm:"foreign_key:KejadianKeselamatanID;references:IdKejadianKeselamatan"`
	orm.Timestamps
}

func (r *ImageKeselamatan) TableName() string {
	return "public.image_keselamatan"
}
