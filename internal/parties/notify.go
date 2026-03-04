package parties

import (
	"fmt"

	"github.com/brunojuliao/go-clappie/internal/tmux"
)

// NotifyStateChange sends a state change message to a participant's pane.
func NotifyStateChange(participant Participant, key, value string) error {
	if participant.PaneID == "" {
		return nil
	}
	msg := fmt.Sprintf("[clappie] State changed → %s = %s", key, value)
	return tmux.SendKeysLiteral(participant.PaneID, msg+"\n")
}

// NotifyAll sends a message to all participants.
func NotifyAll(participants []Participant, message string) {
	for _, p := range participants {
		if p.PaneID != "" {
			tmux.SendKeysLiteral(p.PaneID, message+"\n")
		}
	}
}
