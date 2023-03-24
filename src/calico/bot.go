package calico

import (
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jacobmonck/metrics-collection/src/api/db"
	"github.com/jacobmonck/metrics-collection/src/utils"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	Client    bot.Client
	Config    *utils.Config
	GuildSync Sync
}

type Sync struct {
	Synced        bool
	ResyncChanels bool
	MessageEvents utils.Queue
	MemberEvents  utils.Queue
}

func New(config *utils.Config) (*Bot, error) {
	b := &Bot{
		Config: config,
		GuildSync: Sync{
			Synced: false,
		},
	}

	return b, nil
}

func (b *Bot) Setup(listeners ...bot.EventListener) error {
	client, err := disgo.New(
		utils.RequireEnv("DISCORD_TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentsAll),
		),
		bot.WithEventListeners(append([]bot.EventListener{}, listeners...)...),
	)
	if err != nil {
		return err
	}

	b.Client = client

	return nil
}

func (b *Bot) Start() error {
	err := b.Client.OpenGateway(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) GetGuild() (discord.Guild, error) {
	guild, err := b.Client.Rest().GetGuild(snowflake.ID(b.Config.GuildID), false)
	return guild.Guild, err
}

func (b *Bot) ReplayEvents() {
	totalEvents := len(b.GuildSync.MessageEvents.Items) + len(b.GuildSync.MemberEvents.Items)

	if b.GuildSync.ResyncChanels {
		guild, err := b.GetGuild()
		if err != nil {
			logrus.WithError(err).Error("Failed to fetch guild for event replays.")
			return
		}

		logrus.Info("Resynchronizing channels since an event was sent during startup synchronization.")
		b.SyncChannels(guild)
	}

	for {
		queuedEvent := b.GuildSync.MemberEvents.Pop()
		switch event := queuedEvent.(type) {
		case *events.GuildMemberJoin:
			db.UpdateMember(event.Member, true)
			continue
		case *events.GuildMemberUpdate:
			db.UpdateMember(event.Member, true)
			continue
		case *events.GuildMemberLeave:
			db.UpdateMember(event.Member, false)
			continue
		case nil:
			break
		}
		break
	}

	logrus.Trace("Finished replaying member events.")

	for {
		queuedEvent := b.GuildSync.MessageEvents.Pop()
		switch event := queuedEvent.(type) {
		case *events.GuildMessageCreate:
			db.CreateMessage(event.Message)
			continue
		case *events.GuildMessageDelete:
			db.MarkMessageDeleted(event.MessageID)
			continue
		case nil:
			break
		}
		break
	}

	logrus.Trace("Finished replaying message events.")

	if totalEvents == 0 {
		return
	}

	logrus.Infof("Replayed a total of %d events.", totalEvents)
}
