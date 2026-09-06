package service

import (
	"strings"
	"api-students/app/model"
)

// ValidateCreate memeriksa isi payload pembuatan mahasiswa baru.
func ValidateCreate(req model.CreateStudentRequest) map[string]string {
	errs := map[string]string{}
	if strings.TrimSpace(req.NIM) == "" {
		errs["nim"] = "wajib diisi"
	}
	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "wajib diisi"
	}
	return errs
}

// ValidateReplace memeriksa isi payload pembaruan seluruh data (PUT).
func ValidateReplace(req model.ReplaceStudentRequest) map[string]string {
	errs := map[string]string{}
	if strings.TrimSpace(req.NIM) == "" {
		errs["nim"] = "wajib diisi pada PUT"
	}
	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "wajib diisi pada PUT"
	}
	return errs
}

// IsEmptyPatch memastikan user tidak mengirim JSON kosong saat PATCH.
func IsEmptyPatch(req model.PatchStudentRequest) bool {
	return req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil
}

// ApplyPatch menyalin data baru ke data mahasiswa yang sudah ada dan memvalidasinya.
func ApplyPatch(current model.Student, req model.PatchStudentRequest) (model.Student, map[string]string) {
	errs := map[string]string{}

	if req.NIM != nil {
		if strings.TrimSpace(*req.NIM) == "" {
			errs["nim"] = "tidak boleh kosong"
		} else {
			current.NIM = *req.NIM
		}
	}
	
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			errs["name"] = "tidak boleh kosong"
		} else {
			current.Name = *req.Name
		}
	}
	
	if req.Grade != nil {
		current.Grade = *req.Grade
	}
	
	if req.IsActive != nil {
		current.IsActive = *req.IsActive
	}

	return current, errs
}

// CountTotalPages mengkalkulasi paginasi dengan pembulatan ke atas tanpa float.
func CountTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	return (total + limit - 1) / limit
}