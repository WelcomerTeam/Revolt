package revolt

type User struct {
	ID           string           `json:"_id"`
	Username     string           `json:"username"`
	Avatar       *File            `json:"avatar"`
	Relations    []*UserRelations `json:"relations"`
	Badges       int              `json:"badges"`
	Status       *UserStatus      `json:"status"`
	Bot          *UserBot         `json:"bot,omitempty"`
	Relationship string           `json:"relationship"`
	Online       bool             `json:"online"`
	Flags        int              `json:"flags"`
}

type UserRelations struct {
	Status string `json:"status"`
	UserID string `json:"_id"`
}

type File struct {
	ID          string    `json:"_id"`
	Tag         string    `json:"tag"`
	Size        int       `json:"size"`
	Filename    string    `json:"filename"`
	Metadata    *Metadata `json:"metadata"`
	ContentType string    `json:"content_type"`
}

type UserStatus struct {
	CustomStatus string `json:"text"`
	Presence     string `json:"presence"`
}

type Metadata struct {
	Type   string `json:"type"`
	Width  *int   `json:"width,omitempty"`
	Height *int   `json:"height,omitempty"`
}

type UserBot struct {
	Owner string `json:"owner"`
}

type Guild struct {
	ID                 string `json:"_id"`
	Nonce              string `json:"nonce"`
	Owner              string `json:"owner"`
	Name               string `json:"name"`
	Description        string
	Channels           []string             `json:"channels"`
	Categories         []*GuildCategory     `json:"categories"`
	Roles              []*GuildRole         `json:"roles"`
	SystemMessages     *GuildSystemMessages `json:"system_messages"`
	DefaultPermissions []int                `json:"default_permissions"`

	Icon   *File `json:"icon"`
	Banner *File `json:"banner"`
}

type GuildCategory struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Channels []string `json:"channels"`
}

type GuildRole struct {
	Name        string `json:"name"`
	Permissions []int  `json:"permissions"`
	Colour      string `json:"colour"`
	Hoist       bool   `json:"hoist"`
	Rank        int    `json:"rank"`
}

type GuildMember struct {
	ID *GuildMemberIDs `json:"_id"`
}

type GuildMemberIDs struct {
	Server string `json:"server"`
	User   string `json:"user"`
}

type GuildSystemMessages struct {
	UserJoined string `json:"user_joined"`
	UserLeft   string `json:"user_left"`
	UserKicked string `json:"user_kicked"`
	UserBanned string `json:"user_banned"`
}

type Channel struct {
	ID          string `json:"_id"`
	ChannelType string `json:"channel_type"`
	Server      string `json:"server"`
	Nonce       string `json:"nonce"`
	Name        string `json:"name"`
}

type Message struct {
	ID        string `json:"_id"`
	Nonce     string `json:"nonce"`
	ChannelID string `json:"channel"`
	Author    string `json:"author"`

	// Unfortunately "content" can be a string or object so we have to decode that.

	Content  string `json:"-"`
	TargetID string `json:"-"`
	By       string `json:"-"`
	Name     string `json:"-"`

	ContentType string `json:"-"`

	RawContent interface{} `json:"content"`

	Attachments []*File  `json:"attachments"`
	Edited      int      `json:"edited"`
	Mentions    []string `json:"mentions"`
	Replies     []string `json:"replies"`
}

type MessageRequest struct {
	Content     string   `json:"content"`
	Nonce       string   `json:"nonce"`
	Attachments []string `json:"attachments"`
	Replies     []*Reply `json:"replies"`
}

type Reply struct {
	ID      string `json:"id"`
	Mention bool   `json:"mention"`
}

type MessageContent struct {
	Type string `json:"type"`

	// text
	Content *string `json:"content,omitempty"`

	// system messages
	ID *string `json:"id,omitempty"`

	// system messages/channel_description_changed/channel_icon_changed
	By *string `json:"by,omitempty"`

	// channel_renamed
	Name *string `json:"name,omitempty"`
}
