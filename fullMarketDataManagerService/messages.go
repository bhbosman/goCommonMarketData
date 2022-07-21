package fullMarketDataManagerService

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketData"
	"github.com/bhbosman/goCommsDefinitions"
)

type PublishFullMarketData struct {
	PublishInstrument string
	PubSubBag         goCommsDefinitions.IPubSubBag
}

func NewPublishFullMarketData(publishInstrument string, PubSubBag goCommsDefinitions.IPubSubBag) *PublishFullMarketData {
	return &PublishFullMarketData{
		PublishInstrument: publishInstrument,
		PubSubBag:         PubSubBag,
	}
}

type CallbackMessage struct {
	Data           interface{}
	InstrumentName string
	CallBack       func(data interface{}, fullMarketOrderBook fullMarketData.IFullMarketOrderBook)
}

func NewCallbackMessage(
	instrumentName string,
	callBack func(data interface{}, fullMarketOrderBook fullMarketData.IFullMarketOrderBook),
	data interface{},
) *CallbackMessage {
	return &CallbackMessage{
		Data:           data,
		CallBack:       callBack,
		InstrumentName: instrumentName,
	}
}
