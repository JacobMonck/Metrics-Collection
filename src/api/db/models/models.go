package models

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type Category struct {
	ID       snowflake.ID
	GuildID  snowflake.ID
	Name     string
	Channels []Channel
}

type Channel struct {
	ID         snowflake.ID
	Name       string
	CategoryID snowflake.ID
	Staff      bool
	Threads    []Thread
	Messages   []Message
}

type Thread struct {
	ID                  snowflake.ID
	Name                string
	ChannelID           snowflake.ID
	Archived            bool
	ArchivedAt          time.Time
	AutoArchiveDuration uint16
	Locked              bool
	Type                uint16
	Messages            []*Message
}

type User struct {
	ID           snowflake.ID
	GuildID      snowflake.ID
	Username     string
	Nickname     *string
	AvatarHash   *string
	JoinedAt     time.Time
	CreatedAt    time.Time
	InGuild      bool
	PremiumSince *time.Time
	Pending      bool
	Flags        uint16
	Messages     []Message
}

type Message struct {
	ID        snowflake.ID
	ChannelID snowflake.ID
	ThreadID  *snowflake.ID
	UserID    snowflake.ID
	CreatedAt time.Time
	Deleted   bool
}
