package fullMarketDataManagerService

import (
	"github.com/bhbosman/gocommon/services/IDataShutDown"
	"github.com/bhbosman/gocommon/services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"strings"
)

type InstrumentStatus struct {
	Instrument string
	Status     string
}

type InstrumentStatusArray []InstrumentStatus

func (self InstrumentStatusArray) Len() int {
	return len(self)
}

func (self InstrumentStatusArray) Less(i, j int) bool {
	return strings.Compare(self[i].Instrument, self[j].Instrument) < 0
}

func (self InstrumentStatusArray) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

type FmdBookKey struct {
	FeedName   string
	Instrument string
}

type IFmdManager interface {
	ISendMessage.ISendMessage
	ISendMessage.IMultiSendMessage
	GetInstrumentList() ([]InstrumentStatus, error)
	SubscribeFullMarketData(registerName string, item string)
	UnsubscribeFullMarketData(registerName string, item string)
	SubscribeFullMarketDataMulti(registerName string, items ...string)
	UnsubscribeFullMarketDataMulti(registerName string, items ...string)
}

type IFmdManagerService interface {
	IFmdManager
	IFxService.IFxServices
}

type IFmdManagerData interface {
	IFmdManager
	IDataShutDown.IDataShutDown
}
