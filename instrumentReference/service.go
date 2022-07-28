package instrumentReference

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/pubSub"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type service struct {
	parentContext     context.Context
	ctx               context.Context
	cancelFunc        context.CancelFunc
	cmdChannel        chan interface{}
	onData            func() (IInstrumentReferenceData, error)
	Logger            *zap.Logger
	state             IFxService.State
	pubSub            *pubsub.PubSub
	goFunctionCounter GoFunctionCounter.IService
	subscribeChannel  *pubsub.ChannelSubscription
}

func (self *service) GetKrakenProviders() ([]KrakenReferenceData, error) {
	forProvider, err := CallIInstrumentReferenceGetKrakenProviders(self.ctx, self.cmdChannel, true)
	if err != nil {
		return nil, err
	}
	return forProvider.Args0, forProvider.Args1
}

func (self *service) GetLunoProviders() ([]LunoReferenceData, error) {
	forProvider, err := CallIInstrumentReferenceGetLunoProviders(self.ctx, self.cmdChannel, true)
	if err != nil {
		return nil, err
	}
	return forProvider.Args0, forProvider.Args1
}

func (self *service) Send(message interface{}) error {
	send, err := CallIInstrumentReferenceSend(self.ctx, self.cmdChannel, false, message)
	if err != nil {
		return err
	}
	return send.Args0
}

func (self *service) GetReferenceData(instrumentData string) (IReferenceData, bool) {
	referenceData, err := CallIInstrumentReferenceGetReferenceData(self.ctx, self.cmdChannel, true, instrumentData)
	if err != nil {
		return nil, false
	}
	return referenceData.Args0, referenceData.Args1
}

func (self *service) SomeMethod() {

}

func (self *service) OnStart(ctx context.Context) error {
	err := self.start(ctx)
	if err != nil {
		return err
	}
	self.state = IFxService.Started
	return nil
}

func (self *service) OnStop(ctx context.Context) error {
	err := self.shutdown(ctx)
	close(self.cmdChannel)
	self.state = IFxService.Stopped
	return err
}

func (self *service) shutdown(_ context.Context) error {
	self.cancelFunc()
	return pubSub.Unsubscribe("", self.pubSub, self.goFunctionCounter, self.subscribeChannel)
}

func (self *service) start(_ context.Context) error {
	instanceData, err := self.onData()
	if err != nil {
		return err
	}

	err = self.goFunctionCounter.GoRun(
		"Instrument Reference Service",
		func() {
			self.goStart(instanceData)
		},
	)
	if err != nil {
		return err
	}

	_ = self.Send(self.createInstrumentReference())
	return nil
}

func (self *service) goStart(instanceData IInstrumentReferenceData) {
	defer func(cmdChannel <-chan interface{}) {
		//flush
		for range cmdChannel {
		}
	}(self.cmdChannel)

	self.subscribeChannel = pubsub.NewChannelSubscription(32)
	self.pubSub.AddSub(self.subscribeChannel, self.ServiceName())

	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.ctx,
		instanceData,
		[]ChannelHandler.ChannelHandler{
			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(IInstrumentReference); ok {
						return ChannelEventsForIInstrumentReference(unk, message)
					}
					return false, nil
				},
			},
			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(ISendMessage.ISendMessage); ok {
						return true, unk.Send(message)
					}
					return false, nil
				},
			},
			// TODO: add handlers here
		},
		func() int {
			return len(self.cmdChannel) + self.subscribeChannel.Count()
		},
		goCommsDefinitions.CreateTryNextFunc(self.cmdChannel),
	)
loop:
	for {
		select {
		case <-self.ctx.Done():
			err := instanceData.ShutDown()
			if err != nil {
				self.Logger.Error(
					"error on done",
					zap.Error(err))
			}
			break loop
		case event, ok := <-self.cmdChannel:
			if !ok {
				return
			}
			breakLoop, err := channelHandlerCallback(event)
			if err != nil || breakLoop {
				break loop
			}
		case event, ok := <-self.subscribeChannel.Data:
			if !ok {
				return
			}
			breakLoop, err := channelHandlerCallback(event)
			if err != nil || breakLoop {
				break loop
			}
		}
	}
}

func (self *service) State() IFxService.State {
	return self.state
}

func (self service) ServiceName() string {
	return "InstrumentReference"
}

func (self *service) createInstrumentReference() *MarketDataFeedReference {
	return &MarketDataFeedReference{
		KrakenFeeds: []*KrakenReferenceData{
			{
				ConnectionName: "Connection001",
				Provider:       "Kraken",
				Type:           "book",
				Depth:          100,
				Feeds: []*KrakenFeed{
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.USD/ZAR",
						Pair:             "USD/ZAR",
						MappedInstrument: "USD/ZAR",
					},
					{
						//"Connection XBT/USD",
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.XBT/USD",
						Pair:             "XBT/USD",
						MappedInstrument: "XBT/USD",
					},
					{
						//"Connection XBT/EUR",
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.XBT/EUR",
						Pair:             "XBT/EUR",
						MappedInstrument: "XBT/EUR",
					},
					{
						//"Connection XBT/CAD",
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.XBT/CAD",
						Pair:             "XBT/CAD",
						MappedInstrument: "XBT/CAD",
					},
					{
						//"Connection EUR/USD",
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.EUR/USD",
						Pair:             "EUR/USD",
						MappedInstrument: "EUR/USD",
					},
					{
						//"Connection GBP/USD",
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.GBP/USD",
						Pair:             "GBP/USD",
						MappedInstrument: "GBP/USD",
					},
					{
						//"Connection USD/CAD",
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.USD/CAD",
						Pair:             "USD/CAD",
						MappedInstrument: "USD/CAD",
					},
				},
			},
		},
		LunoFeeds: []*LunoReferenceData{
			// 1.Luno: "BCHXBT"
			{
				SystemName:       "Luno.BCHXBT",
				Provider:         "Luno",
				Name:             "BCHXBT",
				MappedInstrument: "BCH/XBT",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 2.Luno: "ETHAUD"
			{
				SystemName:       "Luno.ETHAUD",
				Provider:         "Luno",
				Name:             "ETHAUD",
				MappedInstrument: "ETH/AUD",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 3.Luno: "ETHEUR"
			{
				SystemName:       "Luno.ETHEUR",
				Provider:         "Luno",
				Name:             "ETHEUR",
				MappedInstrument: "ETH/EUR",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 4.Luno: "ETHGBP"
			{
				SystemName:       "Luno.ETHGBP",
				Provider:         "Luno",
				Name:             "ETHGBP",
				MappedInstrument: "ETH/GBP",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 5.Luno: "ETHIDR"
			{
				SystemName:       "Luno.ETHIDR",
				Provider:         "Luno",
				Name:             "ETHIDR",
				MappedInstrument: "ETH/IDR",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 6.Luno: "ETHMYR"
			{
				SystemName:       "Luno.ETHMYR",
				Provider:         "Luno",
				Name:             "ETHMYR",
				MappedInstrument: "ETH/MYR",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 7.Luno: "ETHNGN"
			{
				SystemName:       "Luno.ETHNGN",
				Provider:         "Luno",
				Name:             "ETHNGN",
				MappedInstrument: "ETH/NGN",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 8.Luno: "ETHUSDC"
			{
				SystemName:       "Luno.ETHUSDC",
				Provider:         "Luno",
				Name:             "ETHUSDC",
				MappedInstrument: "ETH/USDC",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 9.Luno: "ETHXBT"
			{
				SystemName:       "Luno.ETHXBT",
				Provider:         "Luno",
				Name:             "ETHXBT",
				MappedInstrument: "ETH/XBT",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 10.Luno: "ETHZAR"
			{
				SystemName:       "Luno.ETHZAR",
				Provider:         "Luno",
				Name:             "ETHZAR",
				MappedInstrument: "ETH/ZAR",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 11.Luno: "LTCMYR"
			{
				SystemName:       "Luno.LTCMYR",
				Provider:         "Luno",
				Name:             "LTCMYR",
				MappedInstrument: "LTC/MYR",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 12.Luno: "LTCNGN"
			{
				SystemName:       "Luno.LTCNGN",
				Provider:         "Luno",
				Name:             "LTCNGN",
				MappedInstrument: "LTC/NGN",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 13.Luno: "LTCXBT"
			{
				SystemName:       "Luno.LTCXBT",
				Provider:         "Luno",
				Name:             "LTCXBT",
				MappedInstrument: "LTC/XBT",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 14.Luno: "LTCZAR"
			{
				SystemName:       "Luno.LTCZAR",
				Provider:         "Luno",
				Name:             "LTCZAR",
				MappedInstrument: "LTC/ZAR",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 15.Luno: "UNIMYR"
			{
				SystemName:       "Luno.UNIMYR",
				Provider:         "Luno",
				Name:             "UNIMYR",
				MappedInstrument: "UNI/MYR",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 16.Luno: "USDCNGN"
			{
				SystemName:       "Luno.USDCNGN",
				Provider:         "Luno",
				Name:             "USDCNGN",
				MappedInstrument: "USDC/NGN",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 17.Luno: "USDCZAR"
			{
				SystemName:       "Luno.USDCZAR",
				Provider:         "Luno",
				Name:             "USDCZAR",
				MappedInstrument: "USDC/ZAR",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 18.Luno: "XBTAUD"
			{
				SystemName:       "Luno.XBTAUD",
				Provider:         "Luno",
				Name:             "XBTAUD",
				MappedInstrument: "XBT/AUD",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 19.Luno: "XBTEUR"
			{
				SystemName:       "Luno.XBTEUR",
				Provider:         "Luno",
				Name:             "XBTEUR",
				MappedInstrument: "XBT/EUR",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 20.Luno: "XBTGBP"
			{
				SystemName:       "Luno.XBTGBP",
				Provider:         "Luno",
				Name:             "XBTGBP",
				MappedInstrument: "XBT/GBP",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 21.Luno: "XBTIDR"
			{
				SystemName:       "Luno.XBTIDR",
				Provider:         "Luno",
				Name:             "XBTIDR",
				MappedInstrument: "XBT/IDR",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 22.Luno: "XBTMYR"
			{
				SystemName:       "Luno.XBTMYR",
				Provider:         "Luno",
				Name:             "XBTMYR",
				MappedInstrument: "XBT/MYR",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 23.Luno: "XBTNGN"
			{
				SystemName:       "Luno.XBTNGN",
				Provider:         "Luno",
				Name:             "XBTNGN",
				MappedInstrument: "XBT/NGN",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 24.Luno: "XBTUGX"
			{
				SystemName:       "Luno.XBTUGX",
				Provider:         "Luno",
				Name:             "XBTUGX",
				MappedInstrument: "XBT/UGX",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 25.Luno: "XBTZAR"
			{
				SystemName:       "Luno.XBTZAR",
				Provider:         "Luno",
				Name:             "XBTZAR",
				MappedInstrument: "XBT/ZAR",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 26.Luno: "XRPMYR"
			{
				SystemName:       "Luno.XRPMYR",
				Provider:         "Luno",
				Name:             "XRPMYR",
				MappedInstrument: "XRP/MYR",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 27.Luno: "XRPNGN"
			{
				SystemName:       "Luno.XRPNGN",
				Provider:         "Luno",
				Name:             "XRPNGN",
				MappedInstrument: "XRP/NGN",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 28.Luno: "XRPZAR"
			{
				SystemName:       "Luno.XRPZAR",
				Provider:         "Luno",
				Name:             "XRPZAR",
				MappedInstrument: "XRP/ZAR",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 29.Luno: "XRPXBT"
			{
				SystemName:       "Luno.XRPXBT",
				Provider:         "Luno",
				Name:             "XRPXBT",
				MappedInstrument: "XRP/XBT",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
			// 30.Luno: "XBTUSDC"
			{
				SystemName:       "Luno.XBTUSDC",
				Provider:         "Luno",
				Name:             "XBTUSDC",
				MappedInstrument: "XBT/USDC",
				ReferenceData: ReferenceData{
					priceDecimals:  6,
					volumeDecimals: 6,
				},
			},
		},
	}
}

func newService(
	parentContext context.Context,
	onData func() (IInstrumentReferenceData, error),
	logger *zap.Logger,
	pubSub *pubsub.PubSub,
	goFunctionCounter GoFunctionCounter.IService,
) (IInstrumentReferenceService, error) {
	localCtx, localCancelFunc := context.WithCancel(parentContext)
	return &service{
		parentContext:     parentContext,
		ctx:               localCtx,
		cancelFunc:        localCancelFunc,
		cmdChannel:        make(chan interface{}, 32),
		onData:            onData,
		Logger:            logger,
		pubSub:            pubSub,
		goFunctionCounter: goFunctionCounter,
	}, nil
}
