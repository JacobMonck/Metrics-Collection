package calico

import (
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/jacobmonck/metrics-collection/src/utils"
)

type Bot struct {
	Client    bot.Client
	Config    *utils.Config
	GuildSync Sync
}

type Sync struct {
	Synced        bool
	MessageEvents utils.Queue
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
