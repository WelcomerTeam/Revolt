package revolt

import "image/color"

type ImageCreateArguments struct {
	// cache = force_cache or filesize >= filesize_limit
	ForceCache    bool      `json:"force_cache"`
	FilesizeLimit int       `json:"filesize_limit"`
	Options       ImageOpts `json:"options"`
}

type ImageOpts struct {
	// Newline split message
	Text string `json:"text"`

	GuildId int64 `json:"guild_id"`

	ImageURL string `json:"image_url"`

	UserId int64  `json:"user_id"`
	Avatar string `json:"avatar"`

	AllowGIF bool `json:"allow_gif"`

	// Which theme to use when generating images
	Theme Theme `json:"layout"`

	// Identifier for background
	Background string `json:"background"`

	// Identifier for font to use (along with Noto)
	Font string `json:"font"`

	// Border applied to entire image. If transparent, there is no border.
	BorderColour    color.RGBA `json:"-"`
	BorderColourHex string     `json:"border_colour"`
	BorderWidth     int        `json:"border_width"`

	// Alignment of left or right (assuming not vertical layout)
	ProfileAlignment ProfileAlignment `json:"profile_alignment"`

	// Text alignment (left, center, right) (top, middle, bottom)
	TextAlignmentX Xalignment `json:"text_alignment_x"`
	TextAlignmentY Yalignment `json:"text_alignment_y"`

	// Include a border around profile pictures. This also fills
	// under the profile.
	ProfileBorderColour    color.RGBA `json:"-"`
	ProfileBorderColourHex string     `json:"profile_border_colour"`
	// Padding applied to profile pictures inside profile border
	ProfileBorderWidth int `json:"profile_border_width"`
	// Type of curving on the profile border (circle, rounded, square)
	ProfileBorderCurve ProfileBorderCurve `json:"profile_border_curve"`

	// Text stroke. If 0, there is no stroke
	TextStroke          int        `json:"text_stroke"`
	TextStrokeColour    color.RGBA `json:"-"`
	TextStrokeColourHex string     `json:"text_stroke_colour"`

	TextColour    color.RGBA `json:"-"`
	TextColourHex string     `json:"text_colour"`
}

type (
	Xalignment         uint8
	Yalignment         uint8
	ProfileAlignment   uint8
	ProfileBorderCurve uint8
	Theme              uint8
)

const (
	AlignLeft Xalignment = iota
	AlignMiddle
	AlignRight
)

const (
	AlignTop Yalignment = iota
	AlignCenter
	AlignBottom
)

const (
	FloatLeft ProfileAlignment = iota
	FloatRight
)

const (
	CurveCircle ProfileBorderCurve = iota
	CurveSoft
	CurveSquare
)

const (
	ThemeRegular Theme = iota
	ThemeBadge
	ThemeVertical
)
