package hooks

import (
	"context"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// RegisterOAuthHooks registers OAuth2-related hooks
func RegisterOAuthHooks(app *pocketbase.PocketBase) {
	app.OnRecordAfterAuthWithOAuth2Request("users").Add(func(e *core.RecordAuthWithOAuth2Event) error {
		return handleOAuth2Avatar(app, e)
	})
}

// RegisterUserAvatarHook registers hook to reset oauth_avatar_url when user manually changes avatar
func RegisterUserAvatarHook(app *pocketbase.PocketBase) {
	app.OnRecordAfterUpdateRequest("users").Add(func(e *core.RecordUpdateEvent) error {
		oldAvatar := e.Record.OriginalCopy().GetString("avatar")
		newAvatar := e.Record.GetString("avatar")

		if oldAvatar != newAvatar && newAvatar != "" {
			e.Record.Set("oauth_avatar_url", "")

			if err := app.Dao().SaveRecord(e.Record); err != nil {
				app.Logger().Error("Failed to reset oauth_avatar_url", "error", err)
				return err
			}
		}

		return nil
	})
}

// handleOAuth2Avatar implements users avatar synchronization
func handleOAuth2Avatar(app *pocketbase.PocketBase, e *core.RecordAuthWithOAuth2Event) error {
	if e.OAuth2User == nil || e.OAuth2User.AvatarUrl == "" {
		return nil
	}

	avatarURL := e.OAuth2User.AvatarUrl
	savedURL := e.Record.GetString("oauth_avatar_url")
	avatar := e.Record.GetString("avatar")

	shouldDownload := e.IsNewRecord ||
		avatar == "" ||
		(savedURL != "" && avatarURL != savedURL)

	if shouldDownload {
		if err := downloadAvatar(app, e.Record, avatarURL); err != nil {
			app.Logger().Warn("Failed to retrieve OAuth2 avatar",
				"user_id", e.Record.Id,
				"error", err)
		}
	}

	return nil
}

// downloadAvatar downloads and saves user avatar from URL
func downloadAvatar(app *pocketbase.PocketBase, record *models.Record, avatarURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	file, err := filesystem.NewFileFromUrl(ctx, avatarURL)
	if err != nil {
		return err
	}

	form := forms.NewRecordUpsert(app, record)

	if err := form.LoadData(map[string]any{
		"oauth_avatar_url": avatarURL,
	}); err != nil {
		return err
	}

	if err := form.AddFiles("avatar", file); err != nil {
		return err
	}

	if err := form.Submit(); err != nil {
		return err
	}

	return nil
}
