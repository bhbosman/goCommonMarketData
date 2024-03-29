package fullMarketDataManagerViewer

import (
	"context"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/pubSub"
	"github.com/bhbosman/gocommon/services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
)

type service struct {
	ctx                           context.Context
	cancelFunc                    context.CancelFunc
	cmdChannel                    chan interface{}
	onData                        func() (IFullMarketDataViewData, error)
	Logger                        *zap.Logger
	state                         IFxService.State
	pubSub                        *pubsub.PubSub
	goFunctionCounter             GoFunctionCounter.IService
	subscribeChannel              *pubsub.NextFuncSubscription
	onSetMarketDataListChange     func(list []fullMarketDataManagerService.InstrumentStatus) bool
	onSetMarketDataInstanceChange func(data *stream.PublishTop5) bool
	FmdManagerService             fullMarketDataManagerService.IFmdManagerService
	FullMarketDataHelper          fullMarketDataHelper.IFullMarketDataHelper
}

func (self *service) UnsubscribeFullMarketData(item string) {
	name := self.FullMarketDataHelper.InstrumentChannelNameForTop5(item)
	self.pubSub.Unsub(self.subscribeChannel, name)
	//
	self.FmdManagerService.UnsubscribeFullMarketData(
		"View",
		item,
	)

	_, err := CallIFullMarketDataViewUnsubscribeFullMarketData(self.ctx, self.cmdChannel, false, item)
	if err != nil {
		return
	}
}

func (self *service) SubscribeFullMarketData(item string) {
	self.FmdManagerService.SubscribeFullMarketData(
		"View",
		item,
	)
	//
	name := self.FullMarketDataHelper.InstrumentChannelNameForTop5(item)
	self.pubSub.AddSub(self.subscribeChannel, name)

	publishFullMarketData := fullMarketDataManagerService.NewPublishFullMarketData(item, self.subscribeChannel)
	_ = self.FmdManagerService.Send(publishFullMarketData)

	_, err := CallIFullMarketDataViewSubscribeFullMarketData(self.ctx, self.cmdChannel, false, item)
	if err != nil {
		return
	}
}

func (self *service) Send(message interface{}) error {
	send, err := CallIFullMarketDataViewSend(self.ctx, self.cmdChannel, false, message)
	if err != nil {
		return err
	}
	return send.Args0
}

func (self *service) SetMarketDataListChange(change func(list []fullMarketDataManagerService.InstrumentStatus) bool) {
	self.onSetMarketDataListChange = change
}

func (self *service) SetMarketDataInstanceChange(change func(data *stream.PublishTop5) bool) {
	self.onSetMarketDataInstanceChange = change
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
	return pubSub.Unsubscribe("ddd", self.pubSub, self.goFunctionCounter, self.subscribeChannel)
}

func (self *service) State() IFxService.State {
	return self.state
}

func (self *service) ServiceName() string {
	return "MarketDataViewService"
}

func (self *service) start(_ context.Context) error {
	instanceData, err := self.onData()
	instanceData.SetMarketDataListChange(self.onSetMarketDataListChange)
	instanceData.SetMarketDataInstanceChange(self.onSetMarketDataInstanceChange)

	if err != nil {
		return err
	}
	return self.goFunctionCounter.GoRun(
		"Full Depth View Go Function",
		func() {
			self.goStart(instanceData)
		},
	)
}

func (self *service) goStart(instanceData IFullMarketDataViewData) {
	instanceData.Start(self.ServiceName())
	self.subscribeChannel = pubsub.NewNextFuncSubscription(goCommsDefinitions.CreateNextFunc(self.cmdChannel))

	self.pubSub.AddSub(self.subscribeChannel, self.ServiceName(), self.FullMarketDataHelper.FullMarketDataServicePublishInstrumentList())
	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.ctx,
		instanceData,
		[]ChannelHandler.ChannelHandler{
			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(IFullMarketDataView); ok {
						return ChannelEventsForIFullMarketDataView(unk, message)
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
		}
	}
	//flush
	for range self.cmdChannel {
	}
}

func newService(
	parentContext context.Context,
	onData func() (IFullMarketDataViewData, error),
	logger *zap.Logger,
	pubSub *pubsub.PubSub,
	goFunctionCounter GoFunctionCounter.IService,
	FmdManagerService fullMarketDataManagerService.IFmdManagerService,
	FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper,
) (IFullMarketDataViewService, error) {
	localCtx, localCancelFunc := context.WithCancel(parentContext)

	result := &service{
		ctx:                  localCtx,
		cancelFunc:           localCancelFunc,
		cmdChannel:           make(chan interface{}, 32),
		onData:               onData,
		Logger:               logger,
		pubSub:               pubSub,
		goFunctionCounter:    goFunctionCounter,
		FmdManagerService:    FmdManagerService,
		FullMarketDataHelper: FullMarketDataHelper,
	}
	return result, nil
}
