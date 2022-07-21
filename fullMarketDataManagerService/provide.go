package fullMarketDataManagerService

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func Provide(proxy bool) fx.Option {
	return fx.Options(
		fx.Provide(
			func(
				params struct {
					fx.In
					PubSub               *pubsub.PubSub `name:"Application"`
					FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper
				},
			) (func() (IFmdManagerData, error), error) {
				return func() (IFmdManagerData, error) {
					return newData(
						params.PubSub,
						params.FullMarketDataHelper,
						proxy,
					)
				}, nil
			},
		),
		fx.Provide(
			func(
				params struct {
					fx.In
					PubSub                 *pubsub.PubSub  `name:"Application"`
					ApplicationContext     context.Context `name:"Application"`
					OnData                 func() (IFmdManagerData, error)
					Logger                 *zap.Logger
					UniqueReferenceService interfaces.IUniqueReferenceService
					UniqueSessionNumber    interfaces.IUniqueSessionNumber
					GoFunctionCounter      GoFunctionCounter.IService
					FullMarketDataHelper   fullMarketDataHelper.IFullMarketDataHelper
				},
			) (IFmdManagerService, error) {
				serviceInstance, err := newService(
					params.ApplicationContext,
					params.OnData,
					params.Logger,
					params.PubSub,
					params.GoFunctionCounter,
					params.FullMarketDataHelper,
				)
				if err != nil {
					return nil, err
				}
				return serviceInstance, nil
			},
		),
		fx.Invoke(
			func(
				params struct {
					fx.In
					Lifecycle         fx.Lifecycle
					FmdManagerService IFmdManagerService
				},
			) error {
				params.Lifecycle.Append(
					fx.Hook{
						OnStart: params.FmdManagerService.OnStart,
						OnStop:  params.FmdManagerService.OnStop,
					},
				)
				return nil
			},
		),
	)
}
