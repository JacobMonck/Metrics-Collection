package listeners

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/jacobmonck/metrics-collection/src/api/db"
	"github.com/jacobmonck/metrics-collection/src/calico"
	"github.com/sirupsen/logrus"
)

func replayMessageEvents(b *calico.Bot) {
	totalEvents := len(b.GuildSync.MessageEvents.Items)
	for {
		item := b.GuildSync.MessageEvents.Pop()
		switch queuedEvent := item.(type) {
		case *events.GuildMessageCreate:
			db.CreateMessage(queuedEvent.Message, false)
			continue
		case *events.GuildMessageDelete:
			db.MarkMessageDeleted(queuedEvent.MessageID)
			continue
		case nil:
			break
		}
		break
	}

	if totalEvents == 0 {
		return
	}
	logrus.Infof("Replayed a total of %d message events.", totalEvents)
}

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

		replayMessageEvents(b)

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

		replayMessageEvents(b)

		db.MarkMessageDeleted(event.MessageID)
	})
}
