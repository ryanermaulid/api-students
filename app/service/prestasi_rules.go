package service

import (
	"api-students/app/model"
)

func ValidateCreatePrestasi(req model.CreatePrestasiRequest) map[string]string {
	errs := make(map[string]string)

	if req.Name == "" {
		errs["nama_prestasi"] = "nama prestasi tidak boleh kosong"
	}
	if req.Juara != "Lokal" && req.Juara != "Nasional" && req.Juara != "Internasional" {
		errs["juara"] = "juara harus Lokal, Nasional, atau Internasional"
	}

	return errs
}

