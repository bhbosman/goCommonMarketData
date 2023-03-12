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
	messageRouter                 messageRouter.IMessageRouter
	onSetMarketDataListChange     func(list []fullMarketDataManagerService.InstrumentStatus) bool
	onSetMarketDataInstanceChange func(data *stream.PublishTop5) bool
	activeItem                    string
	activeData                    *stream.PublishTop5
	instrumentList                []fullMarketDataManagerService.InstrumentStatus
	instrumentListChanged         bool
}

func (self *data) UnsubscribeFullMarketData(_ string) {
	self.activeItem = ""
}

func (self *data) SubscribeFullMarketData(item string) {
	self.activeItem = item
}

func (self *data) Start(_ string) {
	list, err := self.FmdManagerService.GetInstrumentList()
	if err != nil {
		return
	}
	self.instrumentList = list
	self.instrumentListChanged = true
}

func (self *data) SetMarketDataListChange(change func(list []fullMarketDataManagerService.InstrumentStatus) bool) {
	self.onSetMarketDataListChange = change
}

func (self *data) SetMarketDataInstanceChange(change func(data *stream.PublishTop5) bool) {
	self.onSetMarketDataInstanceChange = change
}

func (self *data) handleFullMarketDepthBookInstruments(message *stream2.FullMarketData_InstrumentList_Response) {
	self.instrumentList = make([]fullMarketDataManagerService.InstrumentStatus, len(message.Instruments), len(message.Instruments))
	for i, instrument := range message.Instruments {
		self.instrumentList[i] = fullMarketDataManagerService.InstrumentStatus{
			Instrument: instrument.Instrument,
			Status:     instrument.Status,
		}
	}
	self.instrumentListChanged = true
}

func (self *data) handlePublishInstanceDataFor(_ *publishInstanceDataFor) error {
	return nil
}

func (self *data) Send(message interface{}) error {
	self.messageRouter.Route(message)
	return nil
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

	_ = result.messageRouter.Add(result.handlePublishTop5)
	_ = result.messageRouter.Add(result.handleEmptyQueue)
	_ = result.messageRouter.Add(result.handlePublishInstanceDataFor)
	_ = result.messageRouter.Add(result.handleFullMarketDepthBookInstruments)
	return result, nil
}
