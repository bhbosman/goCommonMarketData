package instrumentReference

import (
	"context"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(
			func(
				params struct {
					fx.In
				},
			) (func() (IInstrumentReferenceData, error), error) {
				return func() (IInstrumentReferenceData, error) {
					return newData()
				}, nil
			},
		),
		fx.Provide(
			func(
				params struct {
					fx.In
					PubSub                 *pubsub.PubSub  `name:"Application"`
					ApplicationContext     context.Context `name:"Application"`
					OnData                 func() (IInstrumentReferenceData, error)
					Lifecycle              fx.Lifecycle
					Logger                 *zap.Logger
					UniqueReferenceService interfaces.IUniqueReferenceService
					UniqueSessionNumber    interfaces.IUniqueSessionNumber
					GoFunctionCounter      GoFunctionCounter.IService
				},
			) (IInstrumentReferenceService, error) {
				serviceInstance, err := newService(
					params.ApplicationContext,
					params.OnData,
					params.Logger,
					params.PubSub,
					params.GoFunctionCounter,
				)
				if err != nil {
					return nil, err
				}
				params.Lifecycle.Append(
					fx.Hook{
						OnStart: serviceInstance.OnStart,
						OnStop:  serviceInstance.OnStop,
					},
				)
				return serviceInstance, nil
			},
		),
	)
}
