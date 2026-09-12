package service

import (
	"context"
	"strings"
	"time"
	"github.com/gofiber/fiber/v2"
	
	"api-students/app/model"      
	"api-students/app/repository" 
	"api-students/helper"         
)

const refreshTokenBytes = 32

type AuthService struct {
	students   repository.StudentRepository 
	tokens     *repository.TokenRepository   
	jwt        *helper.JWTManager
	refreshTTL time.Duration
}

func NewAuthService(students repository.StudentRepository, tokens *repository.TokenRepository, jwtManager *helper.JWTManager, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		students: students, 
		tokens:   tokens, 
		jwt:      jwtManager, 
		refreshTTL: refreshTTL,
	}
}

func (s *AuthService) Register(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c) 
	defer cancel()

	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body harus berupa JSON")
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if errs := ValidateRegister(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errs})
	}

	// Password DI-HASH sebelum menyentuh database
	hashed, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "gagal memproses password")
	}

	// Buat object mahasiswa baru (Sesuaikan field dengan struct Student lu)
	newStudent := model.Student{
		Name:     req.Username, 
		Email:    req.Email,
		Password: hashed,
		Role:     "user", 
	}

_, err = s.students.Create(ctx, newStudent)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "gagal mendaftarkan mahasiswa")
	}

	return helper.Ok(c, "pendaftaran berhasil", nil)
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body tidak valid")
	}

	if errs := ValidateLogin(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errs})
	}

	// Cari user, sesuaikan method ini di repo lu!
	user, err := s.students.FindByUsername(ctx, strings.TrimSpace(req.Username))
	if err != nil {
		// Tetap jalankan hash palsu agar waktu tanggap mirip[cite: 3]
		helper.VerifyDummyPassword(req.Password)
		return helper.Fail(c, fiber.StatusUnauthorized, "username atau password salah")
	}

	if !helper.VerifyPassword(user.Password, req.Password) {
		return helper.Fail(c, fiber.StatusUnauthorized, "username atau password salah")
	}

	pair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "gagal membuat token")
	}

	return helper.Ok(c, "login berhasil", pair)
}

func (s *AuthService) Refresh(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body tidak valid")
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		return helper.Fail(c, fiber.StatusBadRequest, "refresh_token wajib diisi")
	}

	hash := helper.SHA256Hex(req.RefreshToken)
	stored, err := s.tokens.FindActive(ctx, hash)
	if err != nil {
		return helper.Fail(c, fiber.StatusUnauthorized, "refresh token tidak valid/kedaluwarsa")
	}

	user, err := s.students.FindByID(ctx, stored.UserID) // Sesuaikan method lu
	if err != nil {
		return helper.Fail(c, fiber.StatusUnauthorized, "akun tidak dapat dipakai")
	}

	// ROTASI: token lama dicabut[cite: 3]
	if err := s.tokens.Revoke(ctx, hash); err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "gagal memperbarui token")
	}

	pair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "gagal membuat token")
	}

	return helper.Ok(c, "token diperbarui", pair)
}

func (s *AuthService) Logout(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body tidak valid")
	}

	if strings.TrimSpace(req.RefreshToken) != "" {
		_ = s.tokens.Revoke(ctx, helper.SHA256Hex(req.RefreshToken))
	}
	return helper.Ok(c, "logout berhasil", nil)
}

func (s *AuthService) Me(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	authUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "belum terautentikasi")
	}

	user, err := s.students.FindByID(ctx, authUser.UserID) 
	if err != nil {
		return helper.Fail(c, fiber.StatusUnauthorized, "user tidak ditemukan")
	}

	return helper.Ok(c, "profil diambil", user)
}

func (s *AuthService) issueTokenPair(ctx context.Context, user model.Student) (model.TokenPair, error) {
	accessToken, err := s.jwt.GenerateAccess(user)
	if err != nil {
		return model.TokenPair{}, err
	}

	refreshToken, err := helper.RandomToken(refreshTokenBytes)
	if err != nil {
		return model.TokenPair{}, err
	}

	err = s.tokens.Save(ctx, model.RefreshToken{
		UserID:    user.ID,
		TokenHash: helper.SHA256Hex(refreshToken),
		ExpiresAt: time.Now().Add(s.refreshTTL),
	})
	if err != nil {
		return model.TokenPair{}, err
	}

	return model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwt.AccessTTL().Seconds()),
	}, nil
}