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

func (r *VoteDBRepo) SetVoteValue(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, vote string) (uuid.UUID, error) {
	query := `
	INSERT INTO votes (session_id, user_id, value)
	VALUES ($1, $2, $3)
	ON CONFLICT (session_id, user_id) DO UPDATE
	SET value = EXCLUDED.value,
		updated_at = NOW()
	RETURNING id;
	`

	var voteID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, sessionID, userID, vote).Scan(&voteID)
	if err != nil {
		r.log.Debug("error inserting or updating vote", zap.Error(err))
		return uuid.Nil, err
	}

	return voteID, nil
}
