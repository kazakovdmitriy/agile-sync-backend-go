package converter

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/dbmodel"
	"backend_go/internal/model/entitymodel"
	"github.com/google/uuid"
)

func UserDBToUserSession(user dbmodel.User) apimodel.UsersInSession {
	var id uuid.UUID
	if user.ID != "" {
		parsedID, err := uuid.Parse(user.ID)
		if err == nil {
			id = parsedID
		}
	}

	return apimodel.UsersInSession{
		ID:        id,
		Name:      user.Name,
		IsCreator: user.IsCreator,
		IsWatcher: user.IsWatcher,
	}
}

func UserDBToEntity(user *dbmodel.User) *entitymodel.User {
	if user == nil {
		return nil
	}

	// Конвертируем string в uuid.UUID
	var id uuid.UUID
	if user.ID != "" {
		parsedID, err := uuid.Parse(user.ID)
		if err == nil {
			id = parsedID
		}
	}

	entityUser := &entitymodel.User{
		ID:             id,
		Name:           user.Name,
		Email:          user.Email,          // Конвертируем string в *string
		HashedPassword: user.HashedPassword, // Конвертируем string в *string
		IsActive:       user.IsActive,
		IsVerified:     user.IsVerified,
		IsGuest:        user.IsGuest,
		OAuthID:        user.OAuthID,
		AvatarURL:      user.AvatarURL,
		IsCreator:      user.IsCreator,
		IsWatcher:      user.IsWatcher,
		OnSession:      user.OnSession,
		SocketID:       user.SocketID,
		CreatedAt:      &user.CreatedAt,
	}

	// Конвертируем OAuthProvider
	if user.OAuthProvider != nil {
		provider := entitymodel.OAuthProvider(*user.OAuthProvider)
		entityUser.OAuthProvider = &provider
	}

	// Конвертируем UpdatedAt
	if user.UpdatedAt != nil {
		updatedAt := *user.UpdatedAt
		entityUser.UpdatedAt = &updatedAt
	}

	return entityUser
}

func ToUserProfile(u *entitymodel.User) apimodel.UserProfile {
	var oauthProvider *entitymodel.OAuthProvider
	if u.OAuthProvider != nil {
		// Создаём копию, чтобы избежать мутаций по указателю
		p := *u.OAuthProvider
		oauthProvider = &p
	}

	avatarURL := ""
	if u.AvatarURL != nil {
		avatarURL = *u.AvatarURL
	}

	return apimodel.UserProfile{
		Id:            u.ID,
		Name:          u.Name,
		Email:         u.Email,
		IsActive:      u.IsActive,
		IsVerified:    u.IsVerified,
		IsGuest:       u.IsGuest,
		OAuthProvider: oauthProvider,
		AvatarUrl:     avatarURL,
	}
}
