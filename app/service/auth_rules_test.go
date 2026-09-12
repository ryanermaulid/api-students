package service

import (
	"testing"

	"api-students/app/model"
)

func TestValidateRegister_PasswordTerlaluPendek(t *testing.T) {
	req := model.RegisterRequest{
		Name:     "sari",
		Email:    "sari@example.com",
		Password: "abc123",
	}

	errs := ValidateRegister(req)

	if errs["password"] == "" {
		t.Error("seharusnya ada error password karena kurang dari 8 karakter")
	}
}

func TestValidateRegister_PasswordTanpaAngka(t *testing.T) {
	req := model.RegisterRequest{
		Name:     "sari",
		Email:    "sari@example.com",
		Password: "hanyahuruf", 
	}

	errs := ValidateRegister(req)

	if errs["password"] == "" {
		t.Error("seharusnya ada error password karena tidak memuat angka")
	}
}

func TestValidateRegister_PasswordUmumDitolak(t *testing.T) {
	req := model.RegisterRequest{
		Name:     "sari",
		Email:    "sari@example.com",
		Password: "password123", 
	}

	errs := ValidateRegister(req)

	if errs["password"] == "" {
		t.Error("seharusnya ada error password karena termasuk password umum")
	}
}

func TestValidateRegister_DataValidLolos(t *testing.T) {
	req := model.RegisterRequest{
		Name:     "sari",
		NIM:      "123456789",
		Email:    "sari@example.com",
		Password: "rahasia123", 
	}

	errs := ValidateRegister(req)

	if len(errs) != 0 {
		t.Errorf("data valid seharusnya tidak menghasilkan error, dapat: %v", errs)
	}
}