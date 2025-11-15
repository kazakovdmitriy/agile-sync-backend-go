package repository

import (
	"backend_go/internal/model/converter"
	"backend_go/internal/model/dbmodel"
	"backend_go/internal/model/entitymodel"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type UserDBRepo struct {
	db  *sqlx.DB
	log *zap.Logger
}

func NewUserDBRepo(db *sqlx.DB, log *zap.Logger) *UserDBRepo {
	return &UserDBRepo{db: db, log: log}
}

func (repo *UserDBRepo) Create(ctx context.Context, user *entitymodel.User) (*entitymodel.User, error) {
	query := `
        INSERT INTO users (
            name, email, hashed_password, is_active, is_verified,
            is_guest, is_creator, is_watcher, on_session
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at, updated_at`

	// Используем sql.NullTime для обработки возможных NULL значений
	var (
		id        uuid.UUID
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)

	err := repo.db.QueryRowContext(ctx, query,
		user.Name,
		user.Email,
		user.HashedPassword,
		user.IsActive,
		user.IsVerified,
		user.IsGuest,
		user.IsCreator,
		user.IsWatcher,
		user.OnSession,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Преобразуем sql.NullTime в *time.Time
	user.ID = id
	if createdAt.Valid {
		user.CreatedAt = &createdAt.Time
	} else {
		user.CreatedAt = nil
	}
	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	} else {
		user.UpdatedAt = nil
	}

	return user, nil
}

func (repo *UserDBRepo) GetByEmail(ctx context.Context, email string) (*entitymodel.User, error) {

	query := `
	select id, name, is_creator, socket_id, created_at, updated_at, 
	       is_watcher, on_session, email, hashed_password, is_active, is_verified, 
	       oauth_provider, oauth_id, avatar_url, is_guest
    from users
	where email = $1
	`

	var user dbmodel.User
	err := repo.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}

	return converter.UserDBToEntity(&user), nil
}

func (repo *UserDBRepo) GetByID(ctx context.Context, id uuid.UUID) (*entitymodel.User, error) {
	query := `
	select id, name, is_creator, socket_id, created_at, updated_at, 
	       is_watcher, on_session, email, hashed_password, is_active, is_verified, 
	       oauth_provider, oauth_id, avatar_url, is_guest
    from users
	where id = $1
	`

	var user dbmodel.User
	err := repo.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}

	return converter.UserDBToEntity(&user), nil
}
