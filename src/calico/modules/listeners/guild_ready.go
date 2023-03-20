package listeners

import (
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/jacobmonck/metrics-collection/src/api/db"
	"github.com/jacobmonck/metrics-collection/src/calico"
	"github.com/sirupsen/logrus"
)

func GuildReady(b *calico.Bot, event *events.GuildReady) {
	if uint64(event.GuildID) != b.Config.GuildID {
		return
	}

	go syncChannels(event.Client().Rest(), event.Guild)
	go syncMembers(event.Client(), event.Guild)
}

func syncMembers(client bot.Client, guild discord.Guild) {
	logrus.Info("Synchronizing guild members...")

	apiStart := time.Now()
	members, err := client.MemberChunkingManager().RequestMembersWithQuery(guild.ID, "", guild.MemberCount)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch members from the Discord API.")
	}
	apiDuration := time.Since(apiStart)

	logrus.Infof("Fetched %d members from the Discord API in %s.", len(members), apiDuration)

	logrus.Info("Updating members in database...")

	dbStart := time.Now()
	db.BulkUpsertMembers(members)
	dbDuration := time.Since(dbStart)

	logrus.Infof("Synchronized members with the database in %s.", dbDuration)
}

func syncChannels(rest rest.Rest, guild discord.Guild) {
	logrus.Info("Synchronizing guild channels...")

	channels, err := rest.GetGuildChannels(guild.ID)
	if err != nil {
		logrus.Fatalf("Failed to synchronize channels: %s", err)
	}

	var categoryChannels []discord.GuildCategoryChannel
	var textChannels []discord.GuildTextChannel
	var threadChannels []discord.GuildThread

	for _, channel := range channels {
		switch channel.Type() {
		case discord.ChannelTypeGuildCategory:
			categoryChannels = append(categoryChannels, channel.(discord.GuildCategoryChannel))
		case discord.ChannelTypeGuildText:
			textChannels = append(textChannels, channel.(discord.GuildTextChannel))
		case discord.ChannelTypeGuildPublicThread:
			threadChannels = append(threadChannels, channel.(discord.GuildThread))
		}
	}

	db.UpdateChannels(categoryChannels, textChannels, threadChannels)

	logrus.Info("Finished synchronizing guild channels.")
}
