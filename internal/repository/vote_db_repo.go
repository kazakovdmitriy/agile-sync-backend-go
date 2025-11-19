package repository

import (
	"backend_go/internal/model/entitymodel"
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

func (r *VoteDBRepo) GetBySessionsID(ctx context.Context, sessionID uuid.UUID) ([]*entitymodel.Vote, error) {
	query := `
	select id, session_id, user_id, value, created_at, updated_at
	from votes
	where session_id = $1
	`

	rows, err := r.db.QueryxContext(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	votes := make([]*entitymodel.Vote, 0)
	for rows.Next() {
		var vote entitymodel.Vote
		err := rows.StructScan(&vote)
		if err != nil {
			r.log.Debug("error scan vote", zap.Error(err))
			return nil, err
		}
		votes = append(votes, &vote)
	}

	if err := rows.Err(); err != nil {
		r.log.Debug("Err after scan votes", zap.Error(err))
		return nil, err
	}

	return votes, nil
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
