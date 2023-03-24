package db

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jacobmonck/metrics-collection/src/api/db/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

func BulkUpsertMembers(members []discord.Member) {
	userModels := make([]*models.User, len(members))

	for i, member := range members {
		avatarHash := member.Avatar
		if avatarHash == nil {
			avatarHash = member.User.Avatar
		}

		userModel := &models.User{
			ID:           member.User.ID,
			Username:     member.User.Username,
			Nickname:     member.Nick,
			AvatarHash:   avatarHash,
			JoinedAt:     member.JoinedAt,
			CreatedAt:    member.User.CreatedAt(),
			InGuild:      true,
			PremiumSince: member.PremiumSince,
			Pending:      member.Pending,
			Flags:        uint16(member.Flags),
		}
		userModels[i] = userModel
	}

	DB.Model(&models.User{}).Where("in_guild = ?", true).Update("in_guild", false)

	DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(userModels, 1000)
}

func UpdateMember(member discord.Member, inGuild bool) {
	avatarHash := member.Avatar
	if avatarHash == nil {
		avatarHash = member.User.Avatar
	}

	memberModel := &models.User{
		ID:           member.User.ID,
		Username:     member.User.Username,
		Nickname:     member.Nick,
		AvatarHash:   avatarHash,
		JoinedAt:     member.JoinedAt,
		CreatedAt:    member.User.CreatedAt(),
		InGuild:      inGuild,
		PremiumSince: member.PremiumSince,
		Pending:      member.Pending,
		Flags:        uint16(member.Flags),
	}
	DB.Save(memberModel)
}

func UpdateChannels(
	categoryChannels []discord.GuildCategoryChannel,
	textChannels []discord.GuildTextChannel,
	threadChannels []discord.GuildThread,
) {
	for _, category := range categoryChannels {
		categoryModel := &models.Category{
			ID:   category.ID(),
			Name: category.Name(),
		}
		DB.Save(categoryModel)
	}

	for _, channel := range textChannels {
		channelModel := &models.Channel{
			ID:         channel.ID(),
			Name:       channel.Name(),
			CategoryID: *channel.ParentID(),
		}
		DB.Save(channelModel)
	}

	for _, thread := range threadChannels {
		threadModel := &models.Thread{
			ID:                  thread.ID(),
			Name:                thread.Name(),
			ChannelID:           *thread.ParentID(),
			Archived:            thread.ThreadMetadata.Archived,
			ArchivedAt:          thread.ThreadMetadata.ArchiveTimestamp,
			AutoArchiveDuration: uint16(thread.ThreadMetadata.AutoArchiveDuration),
			Locked:              thread.ThreadMetadata.Locked,
			Type:                uint16(thread.Type()),
		}
		DB.Save(threadModel)
	}
}

func CreateMessage(message discord.Message) {
	var threadID *snowflake.ID
	if thread := message.Thread; thread != nil {
		id := thread.ID()
		threadID = &id
	}

	messageModel := &models.Message{
		ID:        message.ID,
		ChannelID: message.ChannelID,
		ThreadID:  threadID,
		UserID:    message.Author.ID,
		CreatedAt: message.CreatedAt,
	}
	DB.Save(messageModel)
}

func MarkMessageDeleted(messageID snowflake.ID) {
	result := DB.Model(&models.Message{}).
		Where("id = ?", messageID).
		Update("deleted", true)
	if result.RowsAffected == 0 {
		logrus.Warningf("Failed to mark message %d as deleted: no matching rows.", messageID)
	}
}
