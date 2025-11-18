package websocketmodel

// События WebSocket
const (
	EventJoinSession      = "join_session"
	EventVote             = "vote"
	EventRevealCards      = "reveal_cards"
	EventResetVotes       = "reset_votes"
	EventSendReaction     = "send_reaction"
	EventRemoveReaction   = "remove_reaction"
	EventToggleWatcher    = "toggle_watcher_mode"
	EventToggleEmoji      = "toggle_emoji"
	EventToggleAutoReveal = "toggle_auto_reveal"
	EventChangeUsername   = "change_username"
	EventKickUser         = "kick_user"

	// События broadcast
	EventSessionState   = "session_state"
	EventUserJoined     = "user_joined"
	EventVoteSubmitted  = "vote_submitted"
	EventCardsRevealed  = "cards_revealed"
	EventVotesReset     = "votes_reset"
	EventReactionSent   = "reaction_sent"
	EventUserUpdated    = "user_updated"
	EventSessionUpdated = "session_updated"
	EventUserKicked     = "user_kicked"
	EventError          = "error"
)
