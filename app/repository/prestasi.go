package repository

import (
	"context"
	"api-students/app/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PrestasiRepository struct {
	pool *pgxpool.Pool
}

func NewPrestasiRepository(pool *pgxpool.Pool) *PrestasiRepository {
	return &PrestasiRepository{pool: pool}
}

func (r *PrestasiRepository) FindByStudent(ctx context.Context, studentID int) ([]model.Prestasi, error) {
	query := `
		SELECT 
			id_prestasi, 
			student_id, 
			nama_prestasi, 
			juara, 
			is_active, 
			created_at 
		FROM prestasi 
		WHERE student_id = $1
	`
	rows, err := r.pool.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.Prestasi, 0)
	for rows.Next() {
		var p model.Prestasi
		if err := rows.Scan(&p.ID, &p.StudentID, &p.Name, &p.Juara, &p.IsActive, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}	

func (r *PrestasiRepository) Insert(ctx context.Context, req model.CreatePrestasiRequest) error {
	query := `INSERT INTO prestasi (student_id, nama_prestasi, juara, is_active) VALUES ($1, $2, $3, $4)`
	
	isActiveVal := true
	if req.IsActive != nil {
		isActiveVal = *req.IsActive
	}

	_, err := r.pool.Exec(ctx, query, req.StudentID, req.Name, req.Juara, isActiveVal)
	return err
}