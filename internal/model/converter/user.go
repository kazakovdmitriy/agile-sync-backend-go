package converter

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/dbmodel"
	"backend_go/internal/model/entitymodel"
	"github.com/google/uuid"
)

func UserEntityToAPI(user *entitymodel.User) *apimodel.UserResponse {
	if user == nil {
		return nil
	}

	// Конвертируем uuid.UUID в string
	id := user.ID.String()

	apiUser := &apimodel.UserResponse{
		ID:         id,
		Name:       user.Name,
		Email:      *user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		IsGuest:    user.IsGuest,
		AvatarURL:  user.AvatarURL,
		IsCreator:  user.IsCreator,
		IsWatcher:  user.IsWatcher,
		OnSession:  user.OnSession,
		SocketID:   user.SocketID,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	// Копируем OAuth поля, если они не nil
	if user.OAuthProvider != nil {
		oauthProvider := apimodel.OAuthProvider(*user.OAuthProvider)
		apiUser.OAuthProvider = &oauthProvider
	}
	if user.OAuthID != nil {
		oauthID := *user.OAuthID
		apiUser.OAuthID = &oauthID
	}

	return apiUser
}

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
		Email:          &user.Email,         // Конвертируем string в *string
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

func UserAPIToEntity(user *apimodel.UserResponse) *entitymodel.User {
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
		ID:         id,
		Name:       user.Name,
		Email:      &user.Email, // Конвертируем string в *string
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		IsGuest:    user.IsGuest,
		AvatarURL:  user.AvatarURL,
		IsCreator:  user.IsCreator,
		IsWatcher:  user.IsWatcher,
		OnSession:  user.OnSession,
		SocketID:   user.SocketID,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	// Копируем OAuth поля, если они не nil
	if user.OAuthProvider != nil {
		provider := entitymodel.OAuthProvider(*user.OAuthProvider)
		entityUser.OAuthProvider = &provider
	}
	if user.OAuthID != nil {
		oauthID := *user.OAuthID
		entityUser.OAuthID = &oauthID
	}

	return entityUser
}

func UserEntityToDB(user *entitymodel.User) *dbmodel.User {
	if user == nil {
		return nil
	}

	// Конвертируем uuid.UUID в string
	id := user.ID.String()

	dbUser := &dbmodel.User{
		ID:             id,
		Name:           user.Name,
		Email:          *user.Email,         // Конвертируем *string в string
		HashedPassword: user.HashedPassword, // Конвертируем *string в string
		IsActive:       user.IsActive,
		IsVerified:     user.IsVerified,
		IsGuest:        user.IsGuest,
		OAuthID:        user.OAuthID,
		AvatarURL:      user.AvatarURL,
		IsCreator:      user.IsCreator,
		IsWatcher:      user.IsWatcher,
		OnSession:      user.OnSession,
		SocketID:       user.SocketID,
	}

	// Конвертируем OAuthProvider
	if user.OAuthProvider != nil {
		provider := dbmodel.OAuthProviderEnum(*user.OAuthProvider)
		dbUser.OAuthProvider = &provider
	}

	// Конвертируем время
	if user.CreatedAt != nil {
		dbUser.CreatedAt = *user.CreatedAt
	}
	if user.UpdatedAt != nil {
		updatedAt := *user.UpdatedAt
		dbUser.UpdatedAt = &updatedAt
	}

	return dbUser
}
