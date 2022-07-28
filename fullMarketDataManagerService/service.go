package fullMarketDataManagerService

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
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
	proxy                bool
	ctx                  context.Context
	cancelFunc           context.CancelFunc
	cmdChannel           chan interface{}
	onData               func() (IFmdManagerData, error)
	logger               *zap.Logger
	state                IFxService.State
	pubSub               *pubsub.PubSub
	goFunctionCounter    GoFunctionCounter.IService
	subscribeChannel     *pubsub.NextFuncSubscription
	fullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper
}

func (self *service) MultiSend(messages ...interface{}) {
	_, err := CallIFmdManagerMultiSend(self.ctx, self.cmdChannel, false, messages...)
	if err != nil {
		return
	}
}

func (self *service) SubscribeFullMarketData(item string) {
	_, err := CallIFmdManagerSubscribeFullMarketData(self.ctx, self.cmdChannel, false, item)
	if err != nil {
		return
	}
}

func (self *service) UnsubscribeFullMarketData(item string) {
	_, err := CallIFmdManagerUnsubscribeFullMarketData(self.ctx, self.cmdChannel, false, item)
	if err != nil {
		return
	}
}

func (self *service) Send(message interface{}) error {
	send, err := CallIFmdManagerSend(self.ctx, self.cmdChannel, false, message)
	if err != nil {
		return err
	}
	return send.Args0
}

func (self *service) GetInstrumentList() ([]InstrumentStatus, error) {
	list, err := CallIFmdManagerGetInstrumentList(self.ctx, self.cmdChannel, true)
	if err != nil {
		return nil, err
	}
	return list.Args0, list.Args1
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
	return pubSub.Unsubscribe("FMD Manager Service", self.pubSub, self.goFunctionCounter, self.subscribeChannel)
}

func (self *service) start(_ context.Context) error {
	instanceData, err := self.onData()
	if err != nil {
		return err
	}

	return self.goFunctionCounter.GoRun(
		"FMD Manager Service",
		func() {
			self.goStart(instanceData)
		},
	)
}

func (self *service) goStart(instanceData IFmdManagerData) {
	self.subscribeChannel = pubsub.NewNextFuncSubscription(goCommsDefinitions.CreateNextFunc(self.cmdChannel))
	self.pubSub.AddSub(self.subscribeChannel, self.fullMarketDataHelper.FullMarketDataServiceInbound())

	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.ctx,
		instanceData,
		[]ChannelHandler.ChannelHandler{
			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(IFmdManager); ok {
						return ChannelEventsForIFmdManager(unk, message)
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
			n := len(self.cmdChannel) + self.subscribeChannel.Count()
			return n
		},
		goCommsDefinitions.CreateTryNextFunc(self.cmdChannel),
	)
loop:
	for {
		select {
		case <-self.ctx.Done():
			err := instanceData.ShutDown()
			if err != nil {
				self.logger.Error(
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
	// flush
	for range self.cmdChannel {
	}
}

func (self *service) State() IFxService.State {
	return self.state
}

func (self service) ServiceName() string {
	return "FmdManager"
}

func newService(
	parentContext context.Context,
	onData func() (IFmdManagerData, error),
	logger *zap.Logger,
	pubSub *pubsub.PubSub,
	goFunctionCounter GoFunctionCounter.IService,
	fullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper,
) (IFmdManagerService, error) {
	localCtx, localCancelFunc := context.WithCancel(parentContext)
	return &service{
		ctx:                  localCtx,
		cancelFunc:           localCancelFunc,
		cmdChannel:           make(chan interface{}, 32),
		onData:               onData,
		logger:               logger,
		pubSub:               pubSub,
		goFunctionCounter:    goFunctionCounter,
		fullMarketDataHelper: fullMarketDataHelper,
	}, nil
}
