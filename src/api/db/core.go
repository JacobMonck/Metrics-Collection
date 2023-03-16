package db

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/jacobmonck/metrics-collection/src/api/db/models"
	"gorm.io/gorm/clause"
)

func BulkUpsertMembers(members []discord.Member) {
	userModels := make([]*models.User, len(members))

	for i, member := range members {
		var avatarHash *string
		if member.Avatar != nil {
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

	Session.Model(&models.User{}).Where("in_guild = ?", true).Update("in_guild", false)

	Session.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(userModels, 1000)
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

		Session.Save(categoryModel)
	}

	for _, channel := range textChannels {
		channelModel := &models.Channel{
			ID:         channel.ID(),
			Name:       channel.Name(),
			CategoryID: *channel.ParentID(),
		}

		Session.Save(channelModel)
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

		Session.Save(threadModel)
	}
}
