package discordgo

import "time"

// MessageCall stores call metadata for private call messages.
type MessageCall struct {
	// The user IDs that participated in the call.
	Participants []string `json:"participants"`

	// The time at which the call ended, if it has ended.
	EndedTimestamp *time.Time `json:"ended_timestamp"`
}
