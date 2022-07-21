package instrumentReference

import (
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
	parentContext     context.Context
	ctx               context.Context
	cancelFunc        context.CancelFunc
	cmdChannel        chan interface{}
	onData            func() (IInstrumentReferenceData, error)
	Logger            *zap.Logger
	state             IFxService.State
	pubSub            *pubsub.PubSub
	goFunctionCounter GoFunctionCounter.IService
	subscribeChannel  chan interface{}
}

func (self *service) Send(message interface{}) error {
	send, err := CallIInstrumentReferenceSend(self.ctx, self.cmdChannel, false, message)
	if err != nil {
		return err
	}
	return send.Args0
}

func (self *service) GetReferenceData(instrumentData string) (*ReferenceData, bool) {
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

	for _, referenceData := range self.createInstrumentReference() {
		self.Send(referenceData)
	}

	return nil
}

func (self *service) goStart(instanceData IInstrumentReferenceData) {
	defer func(cmdChannel <-chan interface{}) {
		//flush
		for range cmdChannel {
		}
	}(self.cmdChannel)

	subscribeChannel := self.pubSub.Sub(self.ServiceName())

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
			return len(self.cmdChannel) + len(self.subscribeChannel)
		},
		goCommsDefinitions.CreateTryNextFunc(self.cmdChannel),
		//func(i interface{}) {
		//	select {
		//	case self.cmdChannel <- i:
		//		break
		//	default:
		//		break
		//	}
		//},
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
		case event, ok := <-subscribeChannel:
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

func (self *service) createInstrumentReference() []*ReferenceData {

	return []*ReferenceData{
		{
			Name:           "XBTZAR",
			PriceDecimals:  0,
			VolumeDecimals: 6,
		},
		{
			Name:           "ETHXBT",
			PriceDecimals:  6,
			VolumeDecimals: 6,
		},
		{
			Name:           "XBTUGX",
			PriceDecimals:  6,
			VolumeDecimals: 6,
		},
		{
			Name:           "XBTEUR",
			PriceDecimals:  6,
			VolumeDecimals: 6,
		},
		{
			Name:           "XBTZMW",
			PriceDecimals:  6,
			VolumeDecimals: 6,
		},
		{
			Name:           "BCHXBT",
			PriceDecimals:  6,
			VolumeDecimals: 6,
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
