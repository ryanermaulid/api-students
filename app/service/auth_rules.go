package service

import (
	"strings"
	"unicode"
	"api-students/app/model" 
)

const minPasswordLength = 8

func ValidateRegister(req model.RegisterRequest) map[string]string {
	errs := make(map[string]string)
	name := strings.TrimSpace(req.Name)

	switch {
	case name == "":
		errs["name"] = "wajib diisi"
	case len(name) < 3:
		errs["name"] = "minimal 3 karakter"
	case !isValidName(name):
		errs["name"] = "hanya boleh huruf, angka, titik, dan garis bawah"
	}

	nim := strings.TrimSpace(req.NIM)
    if nim == "" {
        errs["nim"] = "wajib diisi"
    } else if len(nim) < 3 {
        errs["nim"] = "minimal 3 karakter"
    }

	// Anggap lu punya fungsi isValidEmail dari tugas sebelumnya, kalau nggak ada hapus blok if ini
	if !strings.Contains(req.Email, "@") {
		errs["email"] = "format email tidak valid"
	}

	if msg := checkPasswordStrength(req.Password); msg != "" {
		errs["password"] = msg
	}

	return errs
}

func ValidateLogin(req model.LoginRequest) map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "wajib diisi"
	}
	if req.Password == "" {
		errs["password"] = "wajib diisi"
	}
	return errs
}

func checkPasswordStrength(password string) string {
	if len(password) < minPasswordLength {
		return "minimal 8 karakter"
	}

	var hasLetter, hasDigit bool
	for _, r := range password {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return "harus memuat huruf dan angka"
	}

	weak := map[string]bool{
		"password1": true, "12345678": true, "qwerty123": true,
		"admin123": true, "password123": true,
	}

	if weak[strings.ToLower(password)] {
		return "password terlalu umum"
	}
	return ""
}

func isValidName(name string) bool {
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '.' && r != '_' {
			return false
		}
	}
	return true
}