package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"api-students/app/model"
)

type TokenRepository struct {
	pool *pgxpool.Pool
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{pool: pool}
}

func (r *TokenRepository) Save(ctx context.Context, t model.RefreshToken) error {
	// PENTING: Pake student_id sesuai tabel lu
	query := `INSERT INTO refresh_tokens (student_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, t.UserID, t.TokenHash, t.ExpiresAt)
	if err != nil {
		return fmt.Errorf("menyimpan refresh token: %w", err)
	}
	return nil
}

func (r *TokenRepository) FindActive(ctx context.Context, tokenHash string) (model.RefreshToken, error) {
	var t model.RefreshToken
	query := `SELECT id, student_id, token_hash, expires_at, revoked_at, created_at 
	          FROM refresh_tokens 
	          WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()`
	
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.RevokedAt, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RefreshToken{}, errors.New("not found")
		}
		return model.RefreshToken{}, fmt.Errorf("mengambil refresh token: %w", err)
	}
	return t, nil
}

func (r *TokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1 AND revoked_at IS NULL`
	_, err := r.pool.Exec(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("mencabut refresh token: %w", err)
	}
	return nil
}