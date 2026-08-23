package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Penyimpanan data sementara di memori
var students []Student
var nextID = 1

func findStudentIndex(id int) int {
	for i := range students {
		if students[i].ID == id {
			return i
		}
	}
	return -1
}

// POST /api/v1/students
func createStudent(c *fiber.Ctx) error {
	var req CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	// Validasi input dasar
	errs := map[string]string{}
	if req.NIM == "" {
		errs["nim"] = "wajib diisi"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "nilai harus antara 0-100"
	}
	if len(errs) > 0 {
		return failValidation(c, errs) 
	}

	// Cek NIM ganda
	for _, s := range students {
		if s.NIM == req.NIM {
			return fail(c, fiber.StatusConflict, "NIM ganda, sudah terdaftar")
		}
	}

	// Simpan data
	baru := Student{
		ID:       nextID,
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: req.IsActive,
	}
	students = append(students, baru)
	nextID++

	// Mengembalikan 201 created
	return created(c, "mahasiswa berhasil ditambah", baru, "/api/v1/students/"+strconv.Itoa(baru.ID))
}