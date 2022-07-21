package fullMarketDataManagerService

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
	"github.com/bhbosman/gocommon/pubSub"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type service struct {
	ctx                  context.Context
	cancelFunc           context.CancelFunc
	cmdChannel           chan interface{}
	onData               func() (IFmdManagerData, error)
	logger               *zap.Logger
	state                IFxService.State
	pubSub               *pubsub.PubSub
	goFunctionCounter    GoFunctionCounter.IService
	subscribeChannel     chan interface{}
	fullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper
}

func (self *service) GetInstrumentList() ([]string, error) {
	list, err := CallIFmdManagerGetInstrumentList(self.ctx, self.cmdChannel, true)
	if err != nil {
		return nil, err
	}
	return list.Args0, list.Args1
}

func (self *service) InstrumentListChannelName() (string, error) {
	name, err := CallIFmdManagerInstrumentListChannelName(self.ctx, self.cmdChannel, true)
	if err != nil {
		return "", err
	}
	return name.Args0, name.Args1
}

func (self *service) InstrumentChannelName(instrument string) (string, error) {
	name, err := CallIFmdManagerInstrumentChannelName(self.ctx, self.cmdChannel, true, instrument)
	if err != nil {
		return "", err
	}
	return name.Args0, name.Args1
}

func (self *service) AllInstrumentChannelName() (string, error) {
	name, err := CallIFmdManagerAllInstrumentChannelName(self.ctx, self.cmdChannel, true)
	if err != nil {
		return "", err
	}
	return name.Args0, name.Args1
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
	self.subscribeChannel = self.pubSub.Sub(self.fullMarketDataHelper.FullMarketDataServiceInbound())
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
			n := len(self.cmdChannel) + len(self.subscribeChannel)
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
		case event, ok := <-self.subscribeChannel:
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
