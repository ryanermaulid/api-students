package service

import (
	"api-students/app/repository"
	"api-students/app/model"
	"api-students/helper"
	"github.com/gofiber/fiber/v2"
)

type PrestasiService struct {
	repo *repository.PrestasiRepository
}

func NewPrestasiService(repo *repository.PrestasiRepository) *PrestasiService {
	return &PrestasiService{repo: repo}
}

func (s *PrestasiService) ListByStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id mahasiswa tidak valid")
	}

	list, err := s.repo.FindByStudent(ctx, id)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "gagal mengambil data prestasi mahasiswa")
	}

	// Kalau list kosong, tetap kembalikan array kosong, bukan error
	if list == nil {
		list = []model.Prestasi{}
	}

	return helper.Ok(c, "data prestasi mahasiswa berhasil diambil", list)
}

func (s *PrestasiService) Create(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	// 1. Ambil ID Student dari parameter URL
	studentID, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id mahasiswa tidak valid")
	}

	// 2. Parsing JSON body ke struct model
	var req model.CreatePrestasiRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusUnprocessableEntity, "format JSON tidak sesuai")
	}

	// Paksa Set StudentID dari parameter URL, jangan percaya input body
	req.StudentID = studentID

	// 3. Panggil fungsi validasi murni
	if errs := ValidateCreatePrestasi(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validasi input gagal",
			"errors":  errs,
		})
	}

	// 4. Eksekusi ke database lewat repository
	if err := s.repo.Insert(ctx, req); err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "gagal menyimpan data prestasi")
	}

	return helper.Ok(c, "data prestasi berhasil ditambahkan", nil)
}