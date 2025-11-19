package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type VoteDBRepo struct {
	db  *sqlx.DB
	log *zap.Logger
}

func NewVoteDBRepo(db *sqlx.DB, log *zap.Logger) *VoteDBRepo {
	return &VoteDBRepo{
		db:  db,
		log: log,
	}
}

func (r *VoteDBRepo) SetVoteValue(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, vote string) error {
	query := `
	INSERT INTO votes (session_id, user_id, value)
	VALUES ($1, $2, $3)
	ON CONFLICT (session_id, user_id) DO UPDATE
	SET value = EXCLUDED.value,
		updated_at = NOW();
	`

	_, err := r.db.ExecContext(ctx, query, sessionID, userID, vote)
	if err != nil {
		r.log.Debug("error insert vote", zap.Error(err))
		return err
	}

	return nil
}
