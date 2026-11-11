package discordgo

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestContentWithMoreMentionsReplaced(t *testing.T) {
	s := &Session{StateEnabled: true, State: NewState()}

	user := &User{
		ID:       "user",
		Username: "User Name",
	}

	s.State.GuildAdd(&Guild{ID: "guild"})
	s.State.RoleAdd("guild", &Role{
		ID:          "role",
		Name:        "Role Name",
		Mentionable: true,
	})
	s.State.MemberAdd(&Member{
		User:    user,
		Nick:    "User Nick",
		GuildID: "guild",
	})
	s.State.ChannelAdd(&Channel{
		Name:    "Channel Name",
		GuildID: "guild",
		ID:      "channel",
	})
	m := &Message{
		Content:      "<@&role> <@!user> <@user> <#channel>",
		ChannelID:    "channel",
		MentionRoles: []string{"role"},
		Mentions:     []*User{user},
	}
	if result, _ := m.ContentWithMoreMentionsReplaced(s); result != "@Role Name @User Nick @User Name #Channel Name" {
		t.Error(result)
	}
}
func TestGettingEmojisFromMessage(t *testing.T) {
	msg := "test test <:kitty14:811736565172011058> <:kitty4:811736468812595260>"
	m := &Message{
		Content: msg,
	}
	emojis := m.GetCustomEmojis()
	if len(emojis) < 1 {
		t.Error("No emojis found.")
		return
	}

}

func TestMessage_Reference(t *testing.T) {
	m := &Message{
		ID:        "811736565172011001",
		GuildID:   "811736565172011002",
		ChannelID: "811736565172011003",
	}

	ref := m.Reference()

	if ref.Type != 0 {
		t.Error("Default reference type should be 0")
	}

	if ref.MessageID != m.ID {
		t.Error("Message ID should be the same")
	}

	if ref.GuildID != m.GuildID {
		t.Error("Guild ID should be the same")
	}

	if ref.ChannelID != m.ChannelID {
		t.Error("Channel ID should be the same")
	}
}

func TestMessage_Forward(t *testing.T) {
	m := &Message{
		ID:        "811736565172011001",
		GuildID:   "811736565172011002",
		ChannelID: "811736565172011003",
	}

	ref := m.Forward()

	if ref.Type != MessageReferenceTypeForward {
		t.Error("Reference type should be 1 (forward)")
	}

	if ref.MessageID != m.ID {
		t.Error("Message ID should be the same")
	}

	if ref.GuildID != m.GuildID {
		t.Error("Guild ID should be the same")
	}

	if ref.ChannelID != m.ChannelID {
		t.Error("Channel ID should be the same")
	}
}

func TestMessageReference_DefaultTypeIsDefault(t *testing.T) {
	r := MessageReference{}
	if r.Type != MessageReferenceTypeDefault {
		t.Error("Default message type should be MessageReferenceTypeDefault")
	}
}

func TestMessage_UnmarshalCall(t *testing.T) {
	var msg Message
	err := json.Unmarshal([]byte(`{
		"id": "123",
		"channel_id": "456",
		"type": 3,
		"timestamp": "2026-06-15T10:00:00.000000+00:00",
		"call": {
			"participants": ["111", "222"],
			"ended_timestamp": "2026-06-15T10:02:03.456000+00:00"
		}
	}`), &msg)
	if err != nil {
		t.Fatal(err)
	}

	if msg.Type != MessageTypeCall {
		t.Fatalf("Message.Type = %d, want %d", msg.Type, MessageTypeCall)
	}
	if msg.Call == nil {
		t.Fatal("Message.Call is nil")
	}
	if want := []string{"111", "222"}; !reflect.DeepEqual(msg.Call.Participants, want) {
		t.Fatalf("Message.Call.Participants = %v, want %v", msg.Call.Participants, want)
	}

	wantEnded, err := time.Parse(time.RFC3339Nano, "2026-06-15T10:02:03.456000+00:00")
	if err != nil {
		t.Fatal(err)
	}
	if msg.Call.EndedTimestamp == nil {
		t.Fatal("Message.Call.EndedTimestamp is nil")
	}
	if !msg.Call.EndedTimestamp.Equal(wantEnded) {
		t.Fatalf("Message.Call.EndedTimestamp = %s, want %s", msg.Call.EndedTimestamp.Format(time.RFC3339Nano), wantEnded.Format(time.RFC3339Nano))
	}
}

func TestMessage_UnmarshalActiveCall(t *testing.T) {
	var msg Message
	err := json.Unmarshal([]byte(`{
		"id": "123",
		"channel_id": "456",
		"type": 3,
		"timestamp": "2026-06-15T10:00:00.000000+00:00",
		"call": {
			"participants": ["111"],
			"ended_timestamp": null
		}
	}`), &msg)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Call == nil {
		t.Fatal("Message.Call is nil")
	}
	if msg.Call.EndedTimestamp != nil {
		t.Fatalf("Message.Call.EndedTimestamp = %s, want nil", msg.Call.EndedTimestamp.Format(time.RFC3339Nano))
	}
}

func TestState_MessageAddMergesCall(t *testing.T) {
	state := NewState()
	state.MaxMessageCount = 10
	err := state.ChannelAdd(&Channel{
		ID:   "456",
		Type: ChannelTypeDM,
	})
	if err != nil {
		t.Fatal(err)
	}

	started := time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)
	ended := started.Add(2*time.Minute + 3*time.Second)
	err = state.MessageAdd(&Message{
		ID:        "123",
		ChannelID: "456",
		Type:      MessageTypeCall,
		Timestamp: started,
		Author: &User{
			ID:       "111",
			Username: "caller",
		},
		Call: &MessageCall{
			Participants: []string{"111"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = state.MessageAdd(&Message{
		ID:        "123",
		ChannelID: "456",
		Call: &MessageCall{
			Participants:   []string{"111", "222"},
			EndedTimestamp: &ended,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	msg, err := state.Message("456", "123")
	if err != nil {
		t.Fatal(err)
	}
	if msg.Author == nil || msg.Author.ID != "111" {
		t.Fatalf("Message.Author = %#v, want caller preserved", msg.Author)
	}
	if msg.Call == nil {
		t.Fatal("Message.Call is nil")
	}
	if want := []string{"111", "222"}; !reflect.DeepEqual(msg.Call.Participants, want) {
		t.Fatalf("Message.Call.Participants = %v, want %v", msg.Call.Participants, want)
	}
	if msg.Call.EndedTimestamp == nil {
		t.Fatal("Message.Call.EndedTimestamp is nil")
	}
	if !msg.Call.EndedTimestamp.Equal(ended) {
		t.Fatalf("Message.Call.EndedTimestamp = %s, want %s", msg.Call.EndedTimestamp.Format(time.RFC3339Nano), ended.Format(time.RFC3339Nano))
	}
}

func TestState_MessageAddMergesPartialCall(t *testing.T) {
	state := NewState()
	state.MaxMessageCount = 10
	err := state.ChannelAdd(&Channel{
		ID:   "456",
		Type: ChannelTypeDM,
	})
	if err != nil {
		t.Fatal(err)
	}

	started := time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)
	ended := started.Add(2*time.Minute + 3*time.Second)
	err = state.MessageAdd(&Message{
		ID:        "123",
		ChannelID: "456",
		Type:      MessageTypeCall,
		Timestamp: started,
		Call: &MessageCall{
			Participants: []string{"111", "222"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = state.MessageAdd(&Message{
		ID:        "123",
		ChannelID: "456",
		Call: &MessageCall{
			EndedTimestamp: &ended,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	msg, err := state.Message("456", "123")
	if err != nil {
		t.Fatal(err)
	}
	if msg.Call == nil {
		t.Fatal("Message.Call is nil")
	}
	if want := []string{"111", "222"}; !reflect.DeepEqual(msg.Call.Participants, want) {
		t.Fatalf("Message.Call.Participants = %v, want %v", msg.Call.Participants, want)
	}
	if msg.Call.EndedTimestamp == nil {
		t.Fatal("Message.Call.EndedTimestamp is nil")
	}
	if !msg.Call.EndedTimestamp.Equal(ended) {
		t.Fatalf("Message.Call.EndedTimestamp = %s, want %s", msg.Call.EndedTimestamp.Format(time.RFC3339Nano), ended.Format(time.RFC3339Nano))
	}
}
