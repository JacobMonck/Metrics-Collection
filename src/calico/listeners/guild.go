package listeners

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jacobmonck/metrics-collection/src/calico"
	"github.com/sirupsen/logrus"
)

func GuildReady(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.GuildReady) {
		if uint64(event.GuildID) != b.Config.GuildID {
			return
		}

		go func() {
			b.SyncGuild(event.Guild)
			b.ReplayEvents()
		}()
	})
}

func GuildAvailable(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.GuildAvailable) {
		if uint64(event.GuildID) != b.Config.GuildID {
			return
		}

		logrus.Info("Guild is available.")

		go func() {
			b.SyncGuild(event.Guild)
			b.ReplayEvents()
		}()
	})
}

func GuildUnavailable(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.GuildUnavailable) {
		if uint64(event.GuildID) != b.Config.GuildID {
			return
		}

		logrus.Warning("Guild is unavailable.")

		b.GuildSync.Synced = false
	})
}

func guildChannelEvent(b *calico.Bot, event any, guildID snowflake.ID) {
	if uint64(guildID) != b.Config.GuildID {
		return
	}

	logrus.Tracef("%T received.", event)

	if b.GuildSync.Synced == false {
		b.GuildSync.ResyncChanels = true
		return
	}

	guild, err := b.Client.Rest().GetGuild(guildID, false)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch guild for channel event.")
		return
	}

	b.SyncChannels(guild.Guild)
}

func GuildChannelCreate(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.GuildChannelCreate) {
		guildChannelEvent(b, event, event.GuildID)
	})
}

func GuildChannelUpdate(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.GuildChannelUpdate) {
		guildChannelEvent(b, event, event.GuildID)
	})
}

func GuildThreadCreate(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.ThreadCreate) {
		guildChannelEvent(b, event, event.GuildID)
	})
}

func GuildThreadUpdate(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.ThreadUpdate) {
		guildChannelEvent(b, event, event.GuildID)
	})
}

func GuildThreadArchive(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.ThreadShow) {
		guildChannelEvent(b, event, event.GuildID)
	})
}

func GuildThreadUnarchive(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.ThreadHide) {
		guildChannelEvent(b, event, event.GuildID)
	})
}
