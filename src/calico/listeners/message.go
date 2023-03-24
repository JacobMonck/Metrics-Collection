package listeners

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jacobmonck/metrics-collection/src/api/db"
	"github.com/jacobmonck/metrics-collection/src/calico"
	"github.com/sirupsen/logrus"
)

func guildMessageEvent(b *calico.Bot, event any, guildID snowflake.ID) {
	if uint64(guildID) != b.Config.GuildID {
		return
	}
	logrus.Tracef("%T received.", event)

	if b.GuildSync.Synced == false {
		b.GuildSync.MessageEvents.Push(event)
		return
	}

	switch messageEvent := event.(type) {
	case *events.GuildMessageCreate:
		db.CreateMessage(messageEvent.Message)
	case *events.GuildMessageDelete:
		db.MarkMessageDeleted(messageEvent.MessageID)
	}
}

func GuildMessageCreate(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.GuildMessageCreate) {
		guildMessageEvent(b, event, event.GuildID)
	})
}

func GuildMessageDelete(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.GuildMessageDelete) {
		guildMessageEvent(b, event, event.GuildID)
	})
}
