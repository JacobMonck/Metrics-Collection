package listeners

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/jacobmonck/metrics-collection/src/api/db"
	"github.com/jacobmonck/metrics-collection/src/calico"
	"github.com/sirupsen/logrus"
)

func GuildMessageCreate(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.GuildMessageCreate) {
		if uint64(event.GuildID) != b.Config.GuildID {
			return
		}

		if b.GuildSync.Synced == false {
			b.GuildSync.MessageEvents.Push(event)
			logrus.Trace("Adding message create event to queue.")
			return
		}

		db.CreateMessage(event.Message, false)
	})
}

func GuildMessageDelete(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.GuildMessageDelete) {
		if uint64(event.GuildID) != b.Config.GuildID {
			return
		}

		if b.GuildSync.Synced == false {
			b.GuildSync.MessageEvents.Push(event)
			logrus.Trace("Adding message delete event to queue.")
			return
		}

		db.MarkMessageDeleted(event.MessageID)
	})
}
