package fullMarketDataHelper

import (
	"go.uber.org/fx"
)

func Provide() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
				},
			) (IFullMarketDataHelper, error) {
				return NewFullMarketDataHelper()
			},
		},
	)
}
