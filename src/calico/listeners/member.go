package listeners

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/jacobmonck/metrics-collection/src/api/db"
	"github.com/jacobmonck/metrics-collection/src/calico"
	"github.com/sirupsen/logrus"
)

func guildMemberEvent(
	b *calico.Bot,
	event any,
	member discord.Member,
	inGuild bool,
) {
	if uint64(member.GuildID) != b.Config.GuildID {
		return
	}
	logrus.Tracef("%T received.", event)

	if b.GuildSync.Synced == false {
		b.GuildSync.MemberEvents.Push(event)
		return
	}

	db.UpdateMember(member, inGuild)
}

func GuildMemberJoin(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.GuildMemberJoin) {
		guildMemberEvent(b, event, event.Member, true)
	})
}

func GuildMemberUpdate(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.GuildMemberUpdate) {
		guildMemberEvent(b, event, event.Member, true)
	})
}

func GuildMemberLeave(b *calico.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.GuildMemberLeave) {
		guildMemberEvent(b, event, event.Member, false)
	})
}
