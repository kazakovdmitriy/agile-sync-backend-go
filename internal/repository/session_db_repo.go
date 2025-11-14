package repository

import (
	"backend_go/internal/model/entitymodel"
	"context"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type SessionDBRepo struct {
	db  *sqlx.DB
	log *zap.Logger
}

func NewSessionDBRepo(db *sqlx.DB, logger *zap.Logger) *SessionDBRepo {
	return &SessionDBRepo{
		db:  db,
		log: logger,
	}
}

func (r *SessionDBRepo) GetByCreator(ctx context.Context, userId string) ([]*entitymodel.Session, error) {
	query := `
	select id, name, deck_type, cards_revealed, 
	       creator_id, creator_name, created_at, 
	       updated_at, allow_emoji, auto_reveal, created_via
		from sessions
		where creator_id = $1
	`

	rows, err := r.db.QueryxContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*entitymodel.Session, 0)
	for rows.Next() {
		var session entitymodel.Session
		if err := rows.StructScan(&session); err != nil {
			r.log.Debug("error message", zap.Error(err))
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	if err := rows.Err(); err != nil {
		r.log.Debug("error message", zap.Error(err))
		return nil, err
	}

	return sessions, nil
}
