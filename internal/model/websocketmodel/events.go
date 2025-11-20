package websocketmodel

type SocketEvent string
type SocketBroadcast string

// События WebSocket
const (
	EventJoinSession      SocketEvent = "join_session"
	EventVote             SocketEvent = "vote"
	EventRevealCards      SocketEvent = "reveal_cards"
	EventResetVotes       SocketEvent = "reset_votes"
	EventSendReaction     SocketEvent = "send_reaction"
	EventRemoveReaction   SocketEvent = "remove_reaction"
	EventToggleWatcher    SocketEvent = "toggle_watcher_mode"
	EventToggleEmoji      SocketEvent = "toggle_emoji"
	EventToggleAutoReveal SocketEvent = "toggle_auto_reveal"
	EventChangeUsername   SocketEvent = "change_username"
	EventKickUser         SocketEvent = "kick_user"

	// События broadcast
	EventSessionState   SocketBroadcast = "session_state"
	EventUserJoined     SocketBroadcast = "user_joined"
	EventVoteSubmitted  SocketBroadcast = "vote_submitted"
	EventCardsRevealed  SocketBroadcast = "cards_revealed"
	EventVotesReset     SocketBroadcast = "votes_reset"
	EventReactionSent   SocketBroadcast = "reaction_sent"
	EventUserUpdated    SocketBroadcast = "user_updated"
	EventSessionUpdated SocketBroadcast = "state_update"
	EventUserKicked     SocketBroadcast = "user_kicked"
	EventError          SocketBroadcast = "error"
)
