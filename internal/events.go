package revolt

type SentBase struct {
	Type string `json:"type"`
}

type Authenticate struct {
	SentBase

	UserID       *string `json:"user_id,omitempty"`
	SessionToken *string `json:"session_token,omitempty"`

	Token *string `json:"token,omitempty"`
}

type BeginTyping struct {
	SentBase

	Channel string `json:"channel"`
}

type EndTyping struct {
	SentBase

	Channel string `json:"channel"`
}

type Ping struct {
	SentBase

	Time int `json:"time"`
}

type Error struct {
	SentBase

	Error string `json:"error"`
}

type Authenticated struct {
	SentBase
}

type Ready struct {
	SentBase

	Users    []*User        `json:"users"`
	Guilds   []*Guild       `json:"servers"`
	Channels []*Channel     `json:"channels"`
	Members  []*GuildMember `json:"members"`
}

type MessageCreate struct {
	SentBase

	*Message
}

type MessageDelete struct {
	SentBase

	MessageID string `json:"id"`
	ChannelID string `json:"channel"`
}

type Pong struct {
	SentBase

	Time int `json:"time"`
}
type MessageUpdate struct {
	SentBase

	ID      string   `json:"id"`
	Message *Message `json:"data"`
}

type ChannelCreate struct {
	SentBase

	*Channel
}

type ChannelUpdate struct {
	SentBase
	ID      string   `json:"id"`
	Channel *Channel `json:"data"`
	Clear   string   `json:"clear"`
}

type ChannelDelete struct {
	SentBase

	ID string `json:"id"`
}

type ChannelGroupJoin struct {
	SentBase

	ChannelID string `json:"id"`
	UserID    string `json:"user"`
}

type ChannelGroupLeave struct {
	SentBase

	ChannelID string `json:"id"`
	UserID    string `json:"user"`
}

type ChannelStartTyping struct {
	SentBase

	ChannelID string `json:"id"`
	UserID    string `json:"user"`
}

type ChannelStopTyping struct {
	SentBase

	ChannelID string `json:"id"`
	UserID    string `json:"user"`
}

type ChannelAck struct {
	SentBase

	ChannelID string `json:"id"`
	UserID    string `json:"user"`
	MessageID string `json:"message_id"`
}

type ServerUpdate struct {
	SentBase

	*Guild
	Clear string `json:"clear"`
}

type ServerDelete struct {
	SentBase

	GuildID string `json:"id"`
}

type ServerMemberUpdate struct {
	SentBase

	GuildID string       `json:"id"`
	Member  *GuildMember `json:"data"`
	Clear   string       `json:"clear"`
}

type ServerMemberJoin struct {
	SentBase

	GuildID string `json:"id"`
	UserID  string `json:"user"`
}

type ServerMemberLeave struct {
	SentBase

	GuildID string `json:"id"`
	UserID  string `json:"user"`
}

type ServerRoleUpdate struct {
	SentBase

	GuildID string `json:"id"`
	UserID  string `json:"user"`
}

type ServerRoleDelete struct {
	SentBase

	GuildID string `json:"id"`
	RoleID  string `json:"role_id"`
}

type UserUpdate struct {
	SentBase

	UserID string `json:"id"`
	Data   *User  `json:"data"`
	Clear  string `json:"clear"`
}

type UserRelationship struct {
	SentBase

	UserID      string `json:"id"`
	OtherUserID string `json:"user"`
	Status      string `json:"status"`
}
