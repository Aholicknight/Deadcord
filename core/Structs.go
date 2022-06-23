package core

import (
	"time"
)

type Science struct {
	Fingerprint string  `json:"fingerprint"`
	Assignments [][]int `json:"assignments"`
}

type GuildChannels []struct {
	ID                   string        `json:"id"`
	Type                 int           `json:"type"`
	Name                 string        `json:"name"`
	Position             int           `json:"position"`
	ParentID             interface{}   `json:"parent_id"`
	GuildID              string        `json:"guild_id"`
	PermissionOverwrites []interface{} `json:"permission_overwrites"`
	Nsfw                 bool          `json:"nsfw"`
	LastMessageID        string        `json:"last_message_id,omitempty"`
	Topic                interface{}   `json:"topic,omitempty"`
	RateLimitPerUser     int           `json:"rate_limit_per_user,omitempty"`
	Banner               interface{}   `json:"banner,omitempty"`
	Bitrate              int           `json:"bitrate,omitempty"`
	UserLimit            int           `json:"user_limit,omitempty"`
	RtcRegion            interface{}   `json:"rtc_region,omitempty"`
}

type RateLimit struct {
	Code       int     `json:"code"`
	Global     bool    `json:"global"`
	Message    string  `json:"message"`
	RetryAfter float64 `json:"retry_after"`
}

type GuildJoin struct {
	Code      string    `json:"code"`
	Type      int       `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
	Guild     Guild     `json:"guild"`
	Channel   Channel   `json:"channel"`
	Inviter   Inviter   `json:"inviter"`
}

type GuildJoinFail struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}

type Inviter struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	PublicFlags   int    `json:"public_flags"`
}

type Invite struct {
	Code      string    `json:"code"`
	Type      int       `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
	Guild     Guild     `json:"guild"`
	Channel   Channel   `json:"channel"`
	Inviter   Inviter   `json:"inviter"`
}

type Message []struct {
	ID              string        `json:"id"`
	Type            int           `json:"type"`
	Content         string        `json:"content"`
	ChannelID       string        `json:"channel_id"`
	Author          Author        `json:"author"`
	Attachments     []interface{} `json:"attachments"`
	Embeds          []interface{} `json:"embeds"`
	Mentions        []interface{} `json:"mentions"`
	MentionRoles    []interface{} `json:"mention_roles"`
	Pinned          bool          `json:"pinned"`
	MentionEveryone bool          `json:"mention_everyone"`
	Tts             bool          `json:"tts"`
	Timestamp       time.Time     `json:"timestamp"`
	EditedTimestamp interface{}   `json:"edited_timestamp"`
	Flags           int           `json:"flags"`
	Components      []interface{} `json:"components"`
	Reactions       []Reactions   `json:"reactions"`
}

type GuildMessages []struct {
	ID              string        `json:"id"`
	Type            int           `json:"type"`
	Content         string        `json:"content"`
	ChannelID       string        `json:"channel_id"`
	Author          Author        `json:"author"`
	Attachments     []interface{} `json:"attachments"`
	Embeds          []interface{} `json:"embeds"`
	Mentions        []interface{} `json:"mentions"`
	MentionRoles    []interface{} `json:"mention_roles"`
	Pinned          bool          `json:"pinned"`
	MentionEveryone bool          `json:"mention_everyone"`
	Tts             bool          `json:"tts"`
	Timestamp       time.Time     `json:"timestamp"`
	EditedTimestamp interface{}   `json:"edited_timestamp"`
	Flags           int           `json:"flags"`
	Components      []interface{} `json:"components"`
}

type Author struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	PublicFlags   int    `json:"public_flags"`
}

type Emoji struct {
	ID   interface{} `json:"id"`
	Name string      `json:"name"`
}

type Reactions struct {
	Emoji Emoji `json:"emoji"`
	Count int   `json:"count"`
	Me    bool  `json:"me"`
}

type CfbmPayload struct {
	M       string        `json:"m"`
	Results []interface{} `json:"results"`
	Timing  int           `json:"timing"`
	Fp      Fp            `json:"fp"`
}

type E struct {
	R  []int `json:"r"`
	Ar []int `json:"ar"`
	Pr int   `json:"pr"`
	Cd int   `json:"cd"`
	Wb bool  `json:"wb"`
	Wp bool  `json:"wp"`
	Wn bool  `json:"wn"`
	Ch bool  `json:"ch"`
	Ws bool  `json:"ws"`
	Wd bool  `json:"wd"`
}

type Fp struct {
	ID int `json:"id"`
	E  E   `json:"e"`
}

type Guild struct {
	ID                          string        `json:"id"`
	Name                        string        `json:"name"`
	Icon                        interface{}   `json:"icon"`
	Description                 interface{}   `json:"description"`
	Splash                      interface{}   `json:"splash"`
	DiscoverySplash             interface{}   `json:"discovery_splash"`
	Features                    []string      `json:"features"`
	Emojis                      []interface{} `json:"emojis"`
	Stickers                    []interface{} `json:"stickers"`
	Banner                      interface{}   `json:"banner"`
	OwnerID                     string        `json:"owner_id"`
	ApplicationID               interface{}   `json:"application_id"`
	Region                      string        `json:"region"`
	AfkChannelID                interface{}   `json:"afk_channel_id"`
	AfkTimeout                  int           `json:"afk_timeout"`
	SystemChannelID             string        `json:"system_channel_id"`
	WidgetEnabled               bool          `json:"widget_enabled"`
	WidgetChannelID             interface{}   `json:"widget_channel_id"`
	VerificationLevel           int           `json:"verification_level"`
	Roles                       []Roles       `json:"roles"`
	DefaultMessageNotifications int           `json:"default_message_notifications"`
	MfaLevel                    int           `json:"mfa_level"`
	ExplicitContentFilter       int           `json:"explicit_content_filter"`
	MaxPresences                interface{}   `json:"max_presences"`
	MaxMembers                  int           `json:"max_members"`
	MaxVideoChannelUsers        int           `json:"max_video_channel_users"`
	VanityURLCode               interface{}   `json:"vanity_url_code"`
	PremiumTier                 int           `json:"premium_tier"`
	PremiumSubscriptionCount    int           `json:"premium_subscription_count"`
	SystemChannelFlags          int           `json:"system_channel_flags"`
	PreferredLocale             string        `json:"preferred_locale"`
	RulesChannelID              string        `json:"rules_channel_id"`
	PublicUpdatesChannelID      string        `json:"public_updates_channel_id"`
	HubType                     interface{}   `json:"hub_type"`
	PremiumProgressBarEnabled   bool          `json:"premium_progress_bar_enabled"`
	Nsfw                        bool          `json:"nsfw"`
	NsfwLevel                   int           `json:"nsfw_level"`
}

type Tags struct {
	BotID string `json:"bot_id"`
}

type Roles struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Permissions  string      `json:"permissions"`
	Position     int         `json:"position"`
	Color        int         `json:"color"`
	Hoist        bool        `json:"hoist"`
	Managed      bool        `json:"managed"`
	Mentionable  bool        `json:"mentionable"`
	Icon         interface{} `json:"icon"`
	UnicodeEmoji interface{} `json:"unicode_emoji"`
	Tags         Tags        `json:"tags,omitempty"`
}

type MembershipScreening struct {
	Version     time.Time    `json:"version"`
	FormFields  []FormFields `json:"form_fields"`
	Description string       `json:"description"`
}

type FormFields struct {
	FieldType string   `json:"field_type"`
	Label     string   `json:"label"`
	Values    []string `json:"values"`
	Required  bool     `json:"required"`
}

type MembershipStatus struct {
	CreatedAt         time.Time   `json:"created_at"`
	ID                string      `json:"id"`
	RejectionReason   interface{} `json:"rejection_reason"`
	ApplicationStatus string      `json:"application_status"`
	LastSeen          time.Time   `json:"last_seen"`
	GuildID           string      `json:"guild_id"`
	UserID            string      `json:"user_id"`
}

type GithubRelease struct {
	URL             string         `json:"url"`
	AssetsURL       string         `json:"assets_url"`
	UploadURL       string         `json:"upload_url"`
	HTMLURL         string         `json:"html_url"`
	ID              int            `json:"id"`
	Author          GithubAuthor   `json:"author"`
	NodeID          string         `json:"node_id"`
	TagName         string         `json:"tag_name"`
	TargetCommitish string         `json:"target_commitish"`
	Name            string         `json:"name"`
	Draft           bool           `json:"draft"`
	Prerelease      bool           `json:"prerelease"`
	CreatedAt       string         `json:"created_at"`
	PublishedAt     string         `json:"published_at"`
	Assets          []GithubAssets `json:"assets"`
	TarballURL      string         `json:"tarball_url"`
	ZipballURL      string         `json:"zipball_url"`
	Body            string         `json:"body"`
}

type GithubAuthor struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type GithubUploader struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type GithubAssets struct {
	URL                string         `json:"url"`
	ID                 int            `json:"id"`
	NodeID             string         `json:"node_id"`
	Name               string         `json:"name"`
	Label              interface{}    `json:"label"`
	Uploader           GithubUploader `json:"uploader"`
	ContentType        string         `json:"content_type"`
	State              string         `json:"state"`
	Size               int            `json:"size"`
	DownloadCount      int            `json:"download_count"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	BrowserDownloadURL string         `json:"browser_download_url"`
}
