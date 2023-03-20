package calico

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/jacobmonck/metrics-collection/src/calico/modules/listeners"
)

func DiscordEventListener(b *Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event bot.Event) {
		switch e := event.(type) {
		case *events.GuildReady:
			listeners.GuildReady(b, e)
		}
		// more events will be added this is used to queue events during sync.
	})
}
