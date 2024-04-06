package seeders

import (
	"goravel/app/models"

	"github.com/goravel/framework/facades"
)

type KejadianSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *KejadianSeeder) Signature() string {
	return "KejadianSeeder"
}

// Run executes the seeder logic.
func (s *KejadianSeeder) Run() error {
	kejadian := make([]models.Kejadian, 0)

	// Kejadian Keamanan Laut
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000002", JenisPelanggaran: "Penyelundupan BMKT", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000009", JenisPelanggaran: "IUU Fishing", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000010", JenisPelanggaran: "Illegal BBM", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000011", JenisPelanggaran: "People Smuggling", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000012", JenisPelanggaran: "Penyelundupan Narkoba", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000019", JenisPelanggaran: "Illegal Logging", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000022", JenisPelanggaran: "Tanpa Izin/Dokumen", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000023", JenisPelanggaran: "Kerusakan Ekosistem", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000024", JenisPelanggaran: "Pelanggaran Wilayah", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000026", JenisPelanggaran: "Penyelundupan Senjata", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000038", JenisPelanggaran: "Perampokan / Pembajakan", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000039", JenisPelanggaran: "Penyelundupan Barang", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000040", JenisPelanggaran: "Illegal Minning", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000041", JenisPelanggaran: "Penyelundupan Miras", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000042", JenisPelanggaran: "Penangkapan ikan menggunakan alat/bom", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000043", JenisPelanggaran: "Penyelundupan Hewan", KlasifikasiID: 1})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000044", JenisPelanggaran: "Human Trafficking", KlasifikasiID: 1})

	// Kejadian Keselamatan Laut
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000007", JenisPelanggaran: "Tabrakan", KlasifikasiID: 2})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000027", JenisPelanggaran: "Tenggelam", KlasifikasiID: 2})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000029", JenisPelanggaran: "Terbakar", KlasifikasiID: 2})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000030", JenisPelanggaran: "Terapung", KlasifikasiID: 2})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000032", JenisPelanggaran: "Karam", KlasifikasiID: 2})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000033", JenisPelanggaran: "Kandas", KlasifikasiID: 2})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000034", JenisPelanggaran: "Hilang Kontak", KlasifikasiID: 2})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000035", JenisPelanggaran: "Terbalik", KlasifikasiID: 2})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000036", JenisPelanggaran: "Hancur", KlasifikasiID: 2})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000045", JenisPelanggaran: "Bocor", KlasifikasiID: 2})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000046", JenisPelanggaran: "Terdampar", KlasifikasiID: 2})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000047", JenisPelanggaran: "Meledak", KlasifikasiID: 2})
	kejadian = append(kejadian, models.Kejadian{IDTypeKejadian: "TYP-000048", JenisPelanggaran: "Kecelakaan Individu", KlasifikasiID: 2})
	
	return facades.Orm().Query().Create(&kejadian)
}
