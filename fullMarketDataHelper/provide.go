package fullMarketDataHelper

import "go.uber.org/fx"

func Provide() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func() (IFullMarketDataHelper, error) {
				return NewFullMarketDataHelper()
			},
		},
	)
}
