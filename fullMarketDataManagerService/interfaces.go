package fullMarketDataManagerService

import (
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/IFxService"
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
	FeedName   string `protobuf:"bytes,1,opt,name=FeedName,proto3" json:"FeedName,omitempty"`
	Instrument string `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
}

type IFmdManager interface {
	ISendMessage.ISendMessage
	ISendMessage.IMultiSendMessage
	GetInstrumentList() ([]InstrumentStatus, error)
	SubscribeFullMarketData(registerName string, item FmdBookKey)
	UnsubscribeFullMarketData(registerName string, item FmdBookKey)
	SubscribeFullMarketDataMulti(registerName string, items ...FmdBookKey)
	UnsubscribeFullMarketDataMulti(registerName string, items ...FmdBookKey)
}

type IFmdManagerService interface {
	IFmdManager
	IFxService.IFxServices
}

type IFmdManagerData interface {
	IFmdManager
	IDataShutDown.IDataShutDown
}
