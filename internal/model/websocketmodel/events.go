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
	EventKickUserError    SocketEvent = "kick_user_error"
	EventUserKicked       SocketEvent = "user_kicked"

	// События broadcast
	EventSessionUpdated      SocketBroadcast = "state_update"
	EventUserKickedBroadcast SocketBroadcast = "user_kicked_broadcast"
)
