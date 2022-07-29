package instrumentReference

import (
	"fmt"
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

	m := make(map[string]bool)
	reference := self.createInstrumentReference()
	for _, krakenFeed := range reference.KrakenFeeds {
		for _, feed := range krakenFeed.Feeds {
			if _, ok := m[feed.Pair]; ok {
				return fmt.Errorf("duplicate instrument in kraken %s", feed.Pair)
			}
			m[feed.Pair] = true
		}
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
					//ADA/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ADA/EUR",
						Pair:             "ADA/EUR",
						MappedInstrument: "ADA/EUR",
					},
					//ADA/ETH
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ADA/ETH",
						Pair:             "ADA/ETH",
						MappedInstrument: "ADA/ETH",
					},
					//ADA/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ADA/USD",
						Pair:             "ADA/USD",
						MappedInstrument: "ADA/USD",
					},
					//ALGO/ETH
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ALGO/ETH",
						Pair:             "ALGO/ETH",
						MappedInstrument: "ALGO/ETH",
					},
					//ADA/XBT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ADA/XBT",
						Pair:             "ADA/XBT",
						MappedInstrument: "ADA/XBT",
					},
					//ALGO/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ALGO/EUR",
						Pair:             "ALGO/EUR",
						MappedInstrument: "ALGO/EUR",
					},
					//ALGO/XBT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ALGO/XBT",
						Pair:             "ALGO/XBT",
						MappedInstrument: "ALGO/XBT",
					},
					//ALGO/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ALGO/USD",
						Pair:             "ALGO/USD",
						MappedInstrument: "ALGO/USD",
					},
					//ATOM/ETH
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ATOM/ETH",
						Pair:             "ATOM/ETH",
						MappedInstrument: "ATOM/ETH",
					},
					//ATOM/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ATOM/USD",
						Pair:             "ATOM/USD",
						MappedInstrument: "ATOM/USD",
					},
					//ATOM/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ATOM/EUR",
						Pair:             "ATOM/EUR",
						MappedInstrument: "ATOM/EUR",
					},
					//ATOM/XBT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ATOM/XBT",
						Pair:             "ATOM/XBT",
						MappedInstrument: "ATOM/XBT",
					},
					//AUD/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.AUD/USD",
						Pair:             "AUD/USD",
						MappedInstrument: "AUD/USD",
					},
					//AUD/JPY
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.AUD/JPY",
						Pair:             "AUD/JPY",
						MappedInstrument: "AUD/JPY",
					},
					//BAL/ETH
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BAL/ETH",
						Pair:             "BAL/ETH",
						MappedInstrument: "BAL/ETH",
					},
					//BAL/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BAL/USD",
						Pair:             "BAL/USD",
						MappedInstrument: "BAL/USD",
					},

					//BAL/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BAL/EUR",
						Pair:             "BAL/EUR",
						MappedInstrument: "BAL/EUR",
					},
					//BAL/XBT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BAL/XBT",
						Pair:             "BAL/XBT",
						MappedInstrument: "BAL/XBT",
					},
					//BAT/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BAT/EUR",
						Pair:             "BAT/EUR",
						MappedInstrument: "BAT/EUR",
					},
					//BAT/ETH
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BAT/ETH",
						Pair:             "BAT/ETH",
						MappedInstrument: "BAT/ETH",
					},
					//BAT/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BAT/USD",
						Pair:             "BAT/USD",
						MappedInstrument: "BAT/USD",
					},
					//BCH/AUD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BCH/AUD",
						Pair:             "BCH/AUD",
						MappedInstrument: "BCH/AUD",
					},
					//BAT/XBT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BAT/XBT",
						Pair:             "BAT/XBT",
						MappedInstrument: "BAT/XBT",
					},
					//BCH/ETH
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BCH/ETH",
						Pair:             "BCH/ETH",
						MappedInstrument: "BCH/ETH",
					},
					//BCH/GBP
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BCH/GBP",
						Pair:             "BCH/GBP",
						MappedInstrument: "BCH/GBP",
					},
					//BCH/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BCH/EUR",
						Pair:             "BCH/EUR",
						MappedInstrument: "BCH/EUR",
					},
					//BCH/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BCH/USD",
						Pair:             "BCH/USD",
						MappedInstrument: "BCH/USD",
					},
					//BCH/USDT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BCH/USDT",
						Pair:             "BCH/USDT",
						MappedInstrument: "BCH/USDT",
					},
					//BCH/XBT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.BCH/XBT",
						Pair:             "BCH/XBT",
						MappedInstrument: "BCH/XBT",
					},
					//COMP/ETH
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.COMP/ETH",
						Pair:             "COMP/ETH",
						MappedInstrument: "COMP/ETH",
					},
					//COMP/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.COMP/EUR",
						Pair:             "COMP/EUR",
						MappedInstrument: "COMP/EUR",
					},
					//COMP/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.COMP/USD",
						Pair:             "COMP/USD",
						MappedInstrument: "COMP/USD",
					},
					//COMP/XBT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.COMP/XBT",
						Pair:             "COMP/XBT",
						MappedInstrument: "COMP/XBT",
					},
					//CRV/ETH
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.CRV/ETH",
						Pair:             "CRV/ETH",
						MappedInstrument: "CRV/ETH",
					},
					//CRV/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.CRV/EUR",
						Pair:             "CRV/EUR",
						MappedInstrument: "CRV/EUR",
					},
					//CRV/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.CRV/USD",
						Pair:             "CRV/USD",
						MappedInstrument: "CRV/USD",
					},
					//CRV/XBT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.CRV/XBT",
						Pair:             "CRV/XBT",
						MappedInstrument: "CRV/XBT",
					},
					//DAI/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.DAI/EUR",
						Pair:             "DAI/EUR",
						MappedInstrument: "DAI/EUR",
					},
					//DAI/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.DAI/USD",
						Pair:             "DAI/USD",
						MappedInstrument: "DAI/USD",
					},
					//DAI/USDT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.DAI/USDT",
						Pair:             "DAI/USDT",
						MappedInstrument: "DAI/USDT",
					},
					//DASH/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.DASH/EUR",
						Pair:             "DASH/EUR",
						MappedInstrument: "DASH/EUR",
					},
					//DASH/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.DASH/USD",
						Pair:             "DASH/USD",
						MappedInstrument: "DASH/USD",
					},
					//DASH/XBT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.DASH/XBT",
						Pair:             "DASH/XBT",
						MappedInstrument: "DASH/XBT",
					},
					//DOT/ETH
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.DOT/ETH",
						Pair:             "DOT/ETH",
						MappedInstrument: "DOT/ETH",
					},
					//DOT/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.DOT/EUR",
						Pair:             "DOT/EUR",
						MappedInstrument: "DOT/EUR",
					},
					//DOT/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.DOT/USD",
						Pair:             "DOT/USD",
						MappedInstrument: "DOT/USD",
					},
					//DOT/XBT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.DOT/XBT",
						Pair:             "DOT/XBT",
						MappedInstrument: "DOT/XBT",
					},
					//EOS/ETH
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.EOS/ETH",
						Pair:             "EOS/ETH",
						MappedInstrument: "EOS/ETH",
					},
					//EOS/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.EOS/EUR",
						Pair:             "EOS/EUR",
						MappedInstrument: "EOS/EUR",
					},
					//EOS/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.EOS/USD",
						Pair:             "EOS/USD",
						MappedInstrument: "EOS/USD",
					},
					//EOS/XBT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.EOS/XBT",
						Pair:             "EOS/XBT",
						MappedInstrument: "EOS/XBT",
					},
					//ETH/AUD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ETH/AUD",
						Pair:             "ETH/AUD",
						MappedInstrument: "ETH/AUD",
					},
					//ETH/CHF
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ETH/CHF",
						Pair:             "ETH/CHF",
						MappedInstrument: "ETH/CHF",
					},
					//ETH/DAI
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ETH/DAI",
						Pair:             "ETH/DAI",
						MappedInstrument: "ETH/DAI",
					},
					//ETH/USDC
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ETH/USDC",
						Pair:             "ETH/USDC",
						MappedInstrument: "ETH/USDC",
					},
					//ETH/USDT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ETH/USDT",
						Pair:             "ETH/USDT",
						MappedInstrument: "ETH/USDT",
					},
					//EUR/AUD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.EUR/AUD",
						Pair:             "EUR/AUD",
						MappedInstrument: "EUR/AUD",
					},
					//EUR/CAD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.EUR/CAD",
						Pair:             "EUR/CAD",
						MappedInstrument: "EUR/CAD",
					},
					//EUR/CHF
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.EUR/CHF",
						Pair:             "EUR/CHF",
						MappedInstrument: "EUR/CHF",
					},
					//EUR/GBP
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.EUR/GBP",
						Pair:             "EUR/GBP",
						MappedInstrument: "EUR/GBP",
					},
					//EUR/JPY
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.EUR/JPY",
						Pair:             "EUR/JPY",
						MappedInstrument: "EUR/JPY",
					},
					//GNO/ETH
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.GNO/ETH",
						Pair:             "GNO/ETH",
						MappedInstrument: "GNO/ETH",
					},
					//GNO/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.GNO/EUR",
						Pair:             "GNO/EUR",
						MappedInstrument: "GNO/EUR",
					},
					//GNO/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.GNO/USD",
						Pair:             "GNO/USD",
						MappedInstrument: "GNO/USD",
					},
					//GNO/XBT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.GNO/XBT",
						Pair:             "GNO/XBT",
						MappedInstrument: "GNO/XBT",
					},
					//ICX/ETH
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ICX/ETH",
						Pair:             "ICX/ETH",
						MappedInstrument: "ICX/ETH",
					},
					//ICX/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ICX/EUR",
						Pair:             "ICX/EUR",
						MappedInstrument: "ICX/EUR",
					},
					//ICX/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ICX/USD",
						Pair:             "ICX/USD",
						MappedInstrument: "ICX/USD",
					},
					//ICX/XBT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.ICX/XBT",
						Pair:             "ICX/XBT",
						MappedInstrument: "ICX/XBT",
					},
					//KAVA/ETH
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.KAVA/ETH",
						Pair:             "KAVA/ETH",
						MappedInstrument: "KAVA/ETH",
					},
					//KAVA/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.KAVA/EUR",
						Pair:             "KAVA/EUR",
						MappedInstrument: "KAVA/EUR",
					},
					//KAVA/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.KAVA/USD",
						Pair:             "KAVA/USD",
						MappedInstrument: "KAVA/USD",
					},
					//KAVA/XBT
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.KAVA/XBT",
						Pair:             "KAVA/XBT",
						MappedInstrument: "KAVA/XBT",
					},
					//KNC/ETH
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.KNC/ETH",
						Pair:             "KNC/ETH",
						MappedInstrument: "KNC/ETH",
					},
					//KNC/EUR
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.KNC/EUR",
						Pair:             "KNC/EUR",
						MappedInstrument: "KNC/EUR",
					},
					//KNC/USD
					{
						ReferenceData: ReferenceData{
							priceDecimals:  6,
							volumeDecimals: 6,
						},
						SystemName:       "Kraken.KNC/USD",
						Pair:             "KNC/USD",
						MappedInstrument: "KNC/USD",
					},
					//KNC/XBT
					//KSM/ETH
					//KSM/EUR
					//KSM/USD
					//KSM/XBT
					//LINK/ETH
					//LINK/EUR
					//LINK/USD
					//LINK/XBT
					//LSK/ETH
					//LSK/EUR
					//LSK/USD
					//LSK/XBT
					//LTC/AUD
					//LTC/ETH
					//LTC/GBP
					//LTC/USDT
					//NANO/ETH
					//NANO/EUR
					//NANO/USD
					//NANO/XBT
					//OMG/ETH
					//OMG/EUR
					//OMG/USD
					//OMG/XBT
					//OXT/ETH
					//OXT/EUR
					//OXT/USD
					//OXT/XBT
					//PAXG/ETH
					//PAXG/EUR
					//PAXG/USD
					//PAXG/XBT
					//QTUM/ETH
					//QTUM/EUR
					//QTUM/USD
					//QTUM/XBT
					//REPV2/ETH
					//REPV2/EUR
					//REPV2/USD
					//REPV2/XBT
					//SC/ETH
					//SC/EUR
					//SC/USD
					//SC/XBT
					//SNX/ETH
					//SNX/EUR
					//SNX/USD
					//SNX/XBT
					//STORJ/ETH
					//STORJ/EUR
					//STORJ/USD
					//STORJ/XBT
					//TRX/ETH
					//TRX/EUR
					//TRX/USD
					//TRX/XBT
					//USDC/EUR
					//USD/CHF
					//USDC/USD
					//USDC/USDT
					//USDT/AUD
					//USDT/CAD
					//USDT/CHF
					//USDT/EUR
					//USDT/GBP
					//USDT/JPY
					//USDT/USD
					//WAVES/ETH
					//WAVES/EUR
					//WAVES/USD
					//WAVES/XBT
					//XBT/AUD
					//XBT/CHF
					//XBT/DAI
					//XBT/USDC
					//XBT/USDT
					//XDG/EUR
					//XDG/USD
					//ETC/ETH
					//ETC/XBT
					//ETC/EUR
					//ETC/USD
					//ETH/XBT
					//ETH/CAD
					//ETH/EUR
					//ETH/GBP
					//ETH/JPY
					//ETH/USD
					//LTC/XBT
					//LTC/EUR
					//LTC/USD
					//MLN/ETH
					//MLN/XBT
					//MLN/EUR
					//MLN/USD
					//REP/ETH
					//REP/XBT
					//REP/EUR
					//REP/USD
					//XRP/AUD
					//XRP/ETH
					//XRP/GBP
					//XRP/USDT
					//XTZ/ETH
					//XTZ/EUR
					//XTZ/USD
					//XTZ/XBT
					//XBT/CAD
					//XBT/EUR
					//XBT/GBP
					//XBT/JPY
					//XBT/USD
					//XDG/XBT
					//XLM/XBT
					//XLM/EUR
					//XLM/USD
					//XMR/XBT
					//XMR/EUR
					//XMR/USD
					//XRP/XBT
					//XRP/CAD
					//XRP/EUR
					//XRP/JPY
					//XRP/USD
					//ZEC/XBT
					//ZEC/EUR
					//ZEC/USD
					//EUR/USD
					//GBP/USD
					//USD/CAD
					//USD/JPY
					//
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
