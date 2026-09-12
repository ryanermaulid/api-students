package model

import "time"

type Prestasi struct {
	ID        int       `json:"id_prestasi"`
	StudentID int       `json:"student_id"`
	Name      string    `json:"nama_prestasi"`
	Juara     string    `json:"juara"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatePrestasiRequest struct {
	StudentID int    `json:"student_id"`
	Name      string `json:"nama_prestasi"`
	Juara     string `json:"juara"`
	IsActive  *bool  `json:"is_active"`
}

type PatchPrestasiRequest struct {
	Name     *string `json:"nama_prestasi"`
	Juara    *string `json:"juara"`
	IsActive *bool   `json:"is_active"`
}

type ReplacePrestasiRequest struct {
	StudentID int    `json:"student_id"`
	Name      string `json:"nama_prestasi"`
	Juara     string `json:"juara"`
	IsActive  bool   `json:"is_active"`
}
