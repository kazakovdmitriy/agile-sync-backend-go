package websocketmodel

type SocketEvent string
type SocketBroadcast string

// События WebSocket
const (
	EventJoinSession   SocketEvent = "join_session"
	EventVote          SocketEvent = "vote"
	EventRevealCards   SocketEvent = "reveal_cards"
	EventResetVotes    SocketEvent = "reset_votes"
	EventKickUser      SocketEvent = "kick_user"
	EventKickUserError SocketEvent = "kick_user_error"
	EventUserKicked    SocketEvent = "user_kicked"
	EventKicked        SocketEvent = "kicked"

	// События broadcast
	EventSessionUpdated      SocketBroadcast = "state_update"
	EventUserKickedBroadcast SocketBroadcast = "user_kicked_broadcast"
)
