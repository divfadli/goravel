package models

import (
	"github.com/goravel/framework/database/orm"
)

type ImageKeamanan struct {
	IDImageKeamanan    uint8            `json:"id_image_keamanan" gorm:"primary_key" column:"id_image_keamanan"`
	FileImageID        uint8            `json:"file_image_id" gorm:"default:0" column:"file_image_id"`
	FileImage          FileImage        `gorm:"foreignKey:FileImageID;references:IdFileImage"`
	KejadianKeamananID uint8            `json:"kejadian_keamanan_id" gorm:"default:0" column:"kejadian_keamanan_id"`
	KejadianKeamanan   KejadianKeamanan `gorm:"foreignKey:KejadianKeamananID;references:IdKejadianKeamanan"`
	orm.Timestamps
}

func (r *ImageKeamanan) TableName() string {
	return "public.image_keamanan"
}
