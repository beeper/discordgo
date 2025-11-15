package discordgo

import (
	"time"

	"github.com/google/uuid"
)

type HeartbeatSession struct {
	CreatedAt         time.Time `json:"createdAtTimestamp"`
	LastUsedTimestamp time.Time `json:"lastUsedTimestamp"`
	ID                uuid.UUID `json:"uuid"`
	// "version": 1
}

func NewHeartbeatSession() HeartbeatSession {
	now := time.Now()
	return HeartbeatSession{
		CreatedAt:         now,
		LastUsedTimestamp: now,
		ID:                uuid.New(),
	}
}

// BumpLastUsed updates the last used timestamp to the current time.
func (hbs *HeartbeatSession) BumpLastUsed() {
	if hbs == nil {
		return
	}
	hbs.LastUsedTimestamp = time.Now()
}

// IsExpired reports whether the heartbeat session should be discarded in favor
// of a new one.
func (hbs *HeartbeatSession) IsExpired() bool {
	if hbs == nil {
		return true
	}
	return time.Since(hbs.LastUsedTimestamp) >= (time.Minute * 30)
}
