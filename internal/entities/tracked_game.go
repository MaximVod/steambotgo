package entities

import "time"

type TrackedGame struct {
	ID          int64
	GameID      int64 // Steam App ID
	GameName    string
	UserChatID  int64 // Telegram Chat ID
	CreatedAt   time.Time
	LastChecked time.Time
}
