// seeds/seeds.go
package seeds

import (
	"fmt"

	"github.com/mfuadfakhruzzaki/backend-api/config"
	"github.com/mfuadfakhruzzaki/backend-api/models"
	"gorm.io/datatypes"
)

func SeedPackages() {
    var count int64
    config.DB.Model(&models.Package{}).Count(&count)
    if count > 0 {
        fmt.Println("Paket sudah ada, skip seeding.")
        return
    }

    packages := []models.Package{
        {
            Name:       "Internet Gatotkaca Ekstra Youtube",
            Data:       "75 GB",
            Duration:   "30 Hari",
            Price:      102000,
            Details:    datatypes.JSON(`["Utama 34GB", "Kuota Lainnya 41GB", "Prime Video 30 Hari"]`),
            Categories: "Paket Gatotkaca",
        },
        {
            Name:       "Internet Gatotkaca Ekstra Youtube",
            Data:       "81 GB",
            Duration:   "30 Hari",
            Price:      112000,
            Details:    datatypes.JSON(`["Utama 39GB", "Kuota Lainnya 42GB", "Prime Video 30 Hari"]`),
            Categories: "Paket Gatotkaca",
        },
        {
            Name:       "Internet Gatotkaca Ekstra Youtube",
            Data:       "85 GB",
            Duration:   "30 Hari",
            Price:      122000,
            Details:    datatypes.JSON(`["Utama 42GB", "Kuota Lainnya 43GB", "Prime Video 30 Hari"]`),
            Categories: "Paket Gatotkaca",
        },
        {
            Name:       "Internet Sebulan + Norton Security",
            Data:       "5 GB",
            Duration:   "30 Hari",
            Price:      56500,
            Details:    datatypes.JSON(`["Utama 5GB", "Prime Video 30 Hari", "Norton Security 30 Hari"]`),
            Categories: "Sebulan",
        },
        {
            Name:       "Internet Sebulan",
            Data:       "12 GB",
            Duration:   "30 Hari",
            Price:      82000,
            Details:    datatypes.JSON(`["Utama 12GB", "Prime Video 30 Hari"]`),
            Categories: "Sebulan",
        },
        {
            Name:       "Internet Sebulan",
            Data:       "39 GB",
            Duration:   "30 Hari",
            Price:      124000,
            Details:    datatypes.JSON(`["Utama 39GB", "Prime Video 30 Hari", "SMS & Voice TSEL"]`),
            Categories: "Sebulan",
        },
        {
            Name:       "Internet Sebulan",
            Data:       "36 GB",
            Duration:   "30 Hari",
            Price:      203000,
            Details:    datatypes.JSON(`["Utama 36GB", "Prime Video 30 Hari", "WeTV 1 Bulan"]`),
            Categories: "Sebulan",
        },
        {
            Name:       "Internet Sebulan",
            Data:       "57 GB",
            Duration:   "30 Hari",
            Price:      164000,
            Details:    datatypes.JSON(`["Utama 57GB", "Prime Video 30 Hari"]`),
            Categories: "Sebulan",
        },
        {
            Name:       "Paket Bundling Zoom Pro",
            Data:       "5 GB",
            Duration:   "30 Hari",
            Price:      54500,
            Details:    datatypes.JSON(`["Utama 5GB", "Zoom Pro 30 Hari"]`),
            Categories: "Lainnya",
        },
        {
            Name:       "Paket Comfort3.5GB",
            Data:       "3.5 GB",
            Duration:   "30 Hari",
            Price:      55000,
            Details:    datatypes.JSON(`["Kuota Lainnya 3.5GB"]`),
            Categories: "Lainnya",
        },
        {
            Name:       "Internet Sebulan",
            Data:       "7.5 GB",
            Duration:   "30 Hari",
            Price:      66500,
            Details:    datatypes.JSON(`["Utama 7.5GB", "Prime Video 30 Hari", "SMS & Voice TSEL"]`),
            Categories: "Sebulan",
        },
        {
            Name:       "Combo GATOTKACA MAX",
            Data:       "4 GB",
            Duration:   "30 Hari",
            Price:      14000,
            Details:    datatypes.JSON(`["Utama 4GB", "Voice & SMS TSEL"]`),
            Categories: "Paket Gatotkaca",
        },
        {
            Name:       "Combo GATOTKACA MAX",
            Data:       "11 GB",
            Duration:   "30 Hari",
            Price:      32000,
            Details:    datatypes.JSON(`["Utama 11GB", "Voucher / Kupon"]`),
            Categories: "Paket Gatotkaca",
        },
        {
            Name:       "Combo GATOTKACA",
            Data:       "17 GB",
            Duration:   "30 Hari",
            Price:      41000,
            Details:    datatypes.JSON(`["Utama 12GB", "Kuota Lainnya 5GB", "Voice & SMS TSEL"]`),
            Categories: "Paket Gatotkaca",
        },
        {
            Name:       "Combo GATOTKACA",
            Data:       "35 GB",
            Duration:   "30 Hari",
            Price:      90000,
            Details:    datatypes.JSON(`["Utama 20GB", "Kuota Lainnya 15GB"]`),
            Categories: "Paket Gatotkaca",
        },
        {
            Name:       "WOW! Nonton 40GB + Vidio",
            Data:       "40 GB",
            Duration:   "30 Hari",
            Price:      112000,
            Details:    datatypes.JSON(`["Utama 13GB", "Kuota Lainnya 13GB", "Vidio"]`),
            Categories: "Paket WOW",
        },
        {
            Name:       "Internet Sebulan",
            Data:       "200 GB",
            Duration:   "30 Hari",
            Price:      450000,
            Details:    datatypes.JSON(`["Utama 200GB", "Prime Video 30 Hari"]`),
            Categories: "Sebulan",
        },
        {
            Name:       "WOW! Chat 6GB + FITA",
            Data:       "6 GB",
            Duration:   "30 Hari",
            Price:      50700,
            Details:    datatypes.JSON(`["Utama 3GB", "Kuota Lainnya 3GB", "FITA"]`),
            Categories: "Paket WOW",
        },
        {
            Name:       "WOW! Nonton 40GB + FITA",
            Data:       "40 GB",
            Duration:   "30 Hari",
            Price:      103000,
            Details:    datatypes.JSON(`["Utama 26GB", "Kuota Lainnya 13GB", "FITA"]`),
            Categories: "Paket WOW",
        },
        {
            Name:       "WOW! Nonton 130GB + Vidio",
            Data:       "130 GB",
            Duration:   "30 Hari",
            Price:      207000,
            Details:    datatypes.JSON(`["Utama 93GB", "Kuota Lainnya 13GB", "Vidio"]`),
            Categories: "Paket WOW",
        },
        {
            Name:       "WOW! Nonton 40GB + Vidio",
            Data:       "40 GB",
            Duration:   "30 Hari",
            Price:      112000,
            Details:    datatypes.JSON(`["Utama 26GB", "Kuota Lainnya 13GB", "Vidio"]`),
            Categories: "Paket WOW",
        },
        {
            Name:       "Paket Data Cihuy Proteksi",
            Data:       "39 GB",
            Duration:   "30 Hari",
            Price:      129000,
            Details:    datatypes.JSON(`["Utama 31GB", "ALLIANZ 30 Hari"]`),
            Categories: "Lainnya",
        },
        {
            Name:       "Internet WOW!",
            Data:       "36 GB",
            Duration:   "30 Hari",
            Price:      110000,
            Details:    datatypes.JSON(`["Utama 18GB", "Prime Video 30 Hari"]`),
            Categories: "Paket WOW",
        },
        {
            Name:       "Paket Data Cihuy Proteksi",
            Data:       "39 GB",
            Duration:   "30 Hari",
            Price:      107000,
            Details:    datatypes.JSON(`["Utama 31GB", "WeTV & ALLIANZ 30 Hari"]`),
            Categories: "Lainnya",
        },
        {
            Name:       "Internet Sebulan",
            Data:       "5 GB",
            Duration:   "30 Hari",
            Price:      54500,
            Details:    datatypes.JSON(`["Utama 5GB", "Prime Video 30 Hari"]`),
            Categories: "Sebulan",
        },
        {
            Name:       "Internet Sebulan",
            Data:       "36 GB",
            Duration:   "30 Hari",
            Price:      199000,
            Details:    datatypes.JSON(`["Utama 36GB", "Prime Video 30 Hari"]`),
            Categories: "Sebulan",
        },
        {
            Name:       "Combo GATOTKACA + Vidio",
            Data:       "48 GB",
            Duration:   "30 Hari",
            Price:      126000,
            Details:    datatypes.JSON(`["Utama 26GB", "Kuota Lainnya 20GB", "Vidio"]`),
            Categories: "Paket Gatotkaca",
        },
    }

    for _, p := range packages {
        config.DB.Create(&p)
    }

    fmt.Println("Seeding data paket selesai.")
}
