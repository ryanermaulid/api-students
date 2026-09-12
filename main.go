package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-students/app/repository"
	"api-students/app/service"
	"api-students/config"
	"api-students/database"
	"api-students/helper" 
)

func main() {
	config.LoadEnv()
	logger := config.NewLogger()

	pool, err := database.NewPool(context.Background())
	if err != nil {
		logger.Error("gagal terhubung ke database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	// 1. Validasi & Init JWT (Jangan sampai tembus kalau secret kosong)
	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		logger.Error("FATAL: JWT_SECRET di .env kependekan atau kosong. Wajib minimal 32 karakter.")
		os.Exit(1)
	}
	jwtManager := helper.NewJWTManager(jwtSecret, "praktikum-backend", 15*time.Minute)

	// 2. Init Repositories
	studentRepo := repository.NewStudentRepository(pool)
	prestasiRepo := repository.NewPrestasiRepository(pool)
	tokenRepo := repository.NewTokenRepository(pool) 

	// 3. Init Services
	studentService := service.NewStudentService(studentRepo)
	prestasiService := service.NewPrestasiService(prestasiRepo)
	
	// Refresh token diset umur 7 hari
	authService := service.NewAuthService(studentRepo, tokenRepo, jwtManager, 7*24*time.Hour) 

	// 4. Injeksi semua ke App 
	app := config.NewApp(logger, pool, jwtManager, studentService, prestasiService, authService)
	
	port := config.GetEnv("APP_PORT", "3000")

	go func() {
		if err := app.Listen(":" + port); err != nil {
			logger.Error("server berhenti", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()
	logger.Info("server berjalan", slog.String("port", port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("sinyal berhenti diterima, menutup server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("gagal menutup server dengan rapi", slog.String("error", err.Error()))
	}
	logger.Info("server berhenti dengan rapi")
}