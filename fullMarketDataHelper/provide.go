package fullMarketDataHelper

import (
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

func Provide() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					PubSub *pubsub.PubSub `name:"Application"`
				},
			) (IFullMarketDataHelper, error) {
				return NewFullMarketDataHelper(params.PubSub)
			},
		},
	)
}
