package fullMarketDataManagerViewer

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goUi/ui"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/services/interfaces"
	"github.com/cskr/pubsub"
	"github.com/rivo/tview"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						FmdManagerService fullMarketDataManagerService.IFmdManagerService
					},
				) (func() (IFullMarketDataViewData, error), error) {
					return func() (IFullMarketDataViewData, error) {
						return newData(params.FmdManagerService)
					}, nil
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						PubSub                 *pubsub.PubSub  `name:"Application"`
						ApplicationContext     context.Context `name:"Application"`
						Lifecycle              fx.Lifecycle
						OnData                 func() (IFullMarketDataViewData, error)
						Logger                 *zap.Logger
						UniqueReferenceService interfaces.IUniqueReferenceService
						UniqueSessionNumber    interfaces.IUniqueSessionNumber
						GoFunctionCounter      GoFunctionCounter.IService
						FmdManagerService      fullMarketDataManagerService.IFmdManagerService
						FullMarketDataHelper   fullMarketDataHelper.IFullMarketDataHelper
					},
				) (IFullMarketDataViewService, error) {
					result, err := newService(
						params.ApplicationContext,
						params.OnData,
						params.Logger,
						params.PubSub,
						params.GoFunctionCounter,
						params.FmdManagerService,
						params.FullMarketDataHelper,
					)
					if err != nil {
						return nil, err
					}
					params.Lifecycle.Append(
						fx.Hook{
							OnStart: result.OnStart,
							OnStop:  result.OnStop,
						})
					return result, nil
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Group: "RegisteredMainWindowSlides",
				Target: func(
					params struct {
						fx.In
						Service                    IFullMarketDataViewService
						App                        *tview.Application
						InstrumentReferenceService instrumentReference.IInstrumentReferenceService
					},
				) (ui.IPrimitiveCloser, error) {
					slide2, err := newSlide(
						2,
						"MarketData",
						params.App,
						params.InstrumentReferenceService,
						params.Service,
					)
					if err != nil {
						return nil, err
					}
					return slide2, nil
				},
			},
		),
	)
}
