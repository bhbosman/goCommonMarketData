package fullMarketDataManagerViewer

import (
	stream2 "github.com/bhbosman/goCommonMarketData/fullMarketData/stream"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
)

type data struct {
	FmdManagerService             fullMarketDataManagerService.IFmdManagerService
	messageRouter                 *messageRouter.MessageRouter
	onSetMarketDataListChange     func(list []string) bool
	onSetMarketDataInstanceChange func(data *stream.PublishTop5) bool
	activeItem                    string
	activeData                    *stream.PublishTop5
	instrumentList                []string
	instrumentListChanged         bool
}

func (self *data) UnsubscribeFullMarketData(item string) {
	self.activeItem = ""
}

func (self *data) SubscribeFullMarketData(item string) {
	self.activeItem = item
}

func (self *data) Start(serviceName string) {
	list, err := self.FmdManagerService.GetInstrumentList()
	if err != nil {
		return
	}
	self.instrumentList = list
	self.instrumentListChanged = true
}

func (self *data) SetMarketDataListChange(change func(list []string) bool) {
	self.onSetMarketDataListChange = change
}

func (self *data) SetMarketDataInstanceChange(change func(data *stream.PublishTop5) bool) {
	self.onSetMarketDataInstanceChange = change
}

func (self *data) handleFullMarketDepthBookInstruments(message *stream2.FullMarketData_InstrumentList_Response) {
	self.instrumentList = message.Instruments
	self.instrumentListChanged = true
}

func (self *data) handlePublishInstanceDataFor(message *publishInstanceDataFor) error {
	return nil
}

func (self *data) Send(message interface{}) error {
	_, err := self.messageRouter.Route(message)
	return err
}

func (self *data) ShutDown() error {
	return nil
}

func (self *data) handlePublishTop5(msg *stream.PublishTop5) {
	if msg.Instrument == self.activeItem {
		self.activeData = msg
	}
}

func (self *data) handleEmptyQueue(msg *messages.EmptyQueue) {
	if self.instrumentListChanged {
		if self.onSetMarketDataListChange != nil {
			b := self.onSetMarketDataListChange(self.instrumentList)
			if !b {
				msg.ErrorHappen = true
				return
			}
		}
		self.instrumentList = nil
		self.instrumentListChanged = false
	}
	if self.activeData != nil {
		if self.onSetMarketDataInstanceChange != nil {
			b := self.onSetMarketDataInstanceChange(self.activeData)
			if !b {
				msg.ErrorHappen = true
			} else {
				self.activeData = nil
			}
		}
	}
}

func newData(FmdManagerService fullMarketDataManagerService.IFmdManagerService) (IFullMarketDataViewData, error) {
	result := &data{
		FmdManagerService: FmdManagerService,
		messageRouter:     messageRouter.NewMessageRouter(),
	}

	result.messageRouter.Add(result.handlePublishTop5)
	result.messageRouter.Add(result.handleEmptyQueue)
	result.messageRouter.Add(result.handlePublishInstanceDataFor)
	result.messageRouter.Add(result.handleFullMarketDepthBookInstruments)
	return result, nil
}
