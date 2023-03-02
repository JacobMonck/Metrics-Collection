package bot

import (
	"context"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/jacobmonck/metrics-collection/src/utils"
)

func Start() (bot.Client, error) {
	client, err := disgo.New(
		utils.RequireEnv("DISCORD_TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentsAll),
		),
	)
	if err != nil {
		return nil, err
	}

	err = client.OpenGateway(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil
}
