package repository

import (
	"backend_go/internal/model/entitymodel"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
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

func (r *SessionDBRepo) CreateSession(ctx context.Context, session *entitymodel.Session) (*entitymodel.Session, error) {
	query := `
	INSERT INTO sessions (
		name, deck_type, cards_revealed, 
		creator_id, creator_name, created_via
	) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, name, deck_type, cards_revealed, 
			  creator_id, creator_name, created_at, 
			  updated_at, allow_emoji, auto_reveal, created_via`

	var sessionResult entitymodel.Session

	err := r.db.QueryRowxContext(ctx, query,
		session.Name,
		session.DeckType,
		session.CardsRevealed,
		session.CreatorID,
		session.CreatorName,
		session.CreatedVia,
	).StructScan(&sessionResult)
	if err != nil {
		return nil, err
	}

	return &sessionResult, nil
}

func (r *SessionDBRepo) GetByID(ctx context.Context, sessionId string) (*entitymodel.Session, error) {
	query := `
	select id, name, deck_type, cards_revealed, 
	       creator_id, creator_name, created_at, 
	       updated_at, allow_emoji, auto_reveal, created_via
		from sessions
		where id = $1
	`

	var sessionResult entitymodel.Session
	err := r.db.QueryRowxContext(ctx, query, sessionId).StructScan(&sessionResult)
	if err != nil {
		return nil, err
	}
	return &sessionResult, nil
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

func (r *SessionDBRepo) DeleteSession(ctx context.Context, sessionId string) error {
	query := `
	DELETE FROM sessions
	where id = $1
	`

	_, err := r.db.ExecContext(ctx, query, sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionDBRepo) ConnectUserToSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	query := `
		INSERT INTO session_connections (session_id, user_id)
		VALUES (:session_id, :user_id)
	`

	// Используем NamedExec для безопасной вставки
	_, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"session_id": sessionID,
		"user_id":    userID,
	})
	if err != nil {
		r.log.Error("Failed to connect user to session",
			zap.String("user_id", userID.String()),
			zap.String("session_id", sessionID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("db insert failed: %w", err)
	}

	return nil
}

func (r *SessionDBRepo) DisconnectUserFromSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	query := `
		UPDATE session_connections
		SET disconnected_at = NOW()
		WHERE 
			session_id = $1 AND 
			user_id = $2 AND 
			disconnected_at IS NULL
		RETURNING id
	`

	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, query, sessionID, userID).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Debug("No active connection to disconnect",
				zap.String("user_id", userID.String()),
				zap.String("session_id", sessionID.String()),
			)
			return nil
		}

		r.log.Error("Failed to disconnect user from session",
			zap.String("user_id", userID.String()),
			zap.String("session_id", sessionID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("db update failed: %w", err)
	}

	return nil
}
