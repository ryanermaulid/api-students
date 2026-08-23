package main

import (
	"sort"
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

// Helper untuk validasi ID dari path parameter
func paramID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

// Fungsi bantu pencarian kata kunci pada Name atau NIM
func cocokPencarian(s Student, kata string) bool {
	kata = strings.ToLower(kata)
	return strings.Contains(strings.ToLower(s.Name), kata) ||
		strings.Contains(strings.ToLower(s.NIM), kata)
}

// GET /api/v1/students/:id
func getStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif") // Status 400
	}

	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan") // Status 404
	}

	return ok(c, "data mahasiswa ditemukan", students[i]) // Status 200
}

// GET /api/v1/students
func listStudents(c *fiber.Ctx) error {
	q := parseListQuery(c)

	// 1) Saring (Filtering)
	hasil := []Student{}
	for _, s := range students {
		if q.IsActive != nil && s.IsActive != *q.IsActive {
			continue
		}
		if q.MinGrade != nil && s.Grade < *q.MinGrade {
			continue
		}
		if q.Search != "" && !cocokPencarian(s, q.Search) {
			continue
		}
		hasil = append(hasil, s)
	}

	// 2) Urutkan (Sorting dengan Whitelist)
	sort.SliceStable(hasil, func(i, j int) bool {
		var lebihKecil bool
		switch q.Sort {
		case "name":
			lebihKecil = strings.ToLower(hasil[i].Name) < strings.ToLower(hasil[j].Name)
		case "nim":
			lebihKecil = hasil[i].NIM < hasil[j].NIM
		case "grade":
			lebihKecil = hasil[i].Grade < hasil[j].Grade
		default:
			lebihKecil = hasil[i].ID < hasil[j].ID
		}

		if q.Order == "desc" {
			return !lebihKecil
		}
		return lebihKecil
	})

	// 3) pagination
	total := len(hasil)
	totalPages := (total + q.Limit - 1) / q.Limit
	mulai := (q.Page - 1) * q.Limit
	if mulai > total {
		mulai = total
	}
	akhir := mulai + q.Limit
	if akhir > total {
		akhir = total
	}

	return okList(c, "daftar mahasiswa berhasil diambil", hasil[mulai:akhir], &Meta{
		Page:       q.Page,
		Limit:      q.Limit,
		Total:      total,
		TotalPages: totalPages,
	})
}