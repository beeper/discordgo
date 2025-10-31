package discordgo

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// TestLaunchSignatureValidUUID ensures that launch signatures double as valid
// version 4 UUIDs.
func TestLaunchSignatureValidUUID(t *testing.T) {
	sig, err := NewVanillaSignature()
	if err != nil {
		t.Errorf(err.Error())
	}
	uuid.MustParse(sig.String())
}

func TestLaunchSignatureMarshal(t *testing.T) {
	sig, err := NewVanillaSignature()
	if err != nil {
		t.Errorf(err.Error())
	}
	j, err := json.Marshal(sig)
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(j) != fmt.Sprintf("\"%s\"", sig.String()) {
		t.Errorf("launch signature didn't serialize as UUID; got %v", string(j))
	}
}
