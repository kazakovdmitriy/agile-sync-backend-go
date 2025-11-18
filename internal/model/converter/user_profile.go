package converter

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/entitymodel"
)

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
		Email:         *u.Email,
		IsActive:      u.IsActive,
		IsVerified:    u.IsVerified,
		IsGuest:       u.IsGuest,
		OAuthProvider: oauthProvider,
		AvatarUrl:     avatarURL,
	}
}
