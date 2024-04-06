package seeders

import (
	"goravel/app/models"

	"github.com/goravel/framework/facades"
)

type KlasifikasiKejadianSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *KlasifikasiKejadianSeeder) Signature() string {
	return "KlasifikasiKejadianSeeder"
}

// Run executes the seeder logic.
func (s *KlasifikasiKejadianSeeder) Run() error {
	klasifikasi := make([]models.KlasifikasiKejadian, 0)

	// Klasifikasi Kejadian
	klasifikasi = append(klasifikasi, models.KlasifikasiKejadian{NamaKlasifikasi: "Keamanan Laut"})
	klasifikasi = append(klasifikasi, models.KlasifikasiKejadian{NamaKlasifikasi: "Keselamatan Laut"})

	return facades.Orm().Query().Create(&klasifikasi)
}
