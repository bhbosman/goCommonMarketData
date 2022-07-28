package fullMarketDataManagerService

import (
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
)

type InstrumentStatus struct {
	Instrument string
	Status     string
}

type IFmdManager interface {
	ISendMessage.ISendMessage
	ISendMessage.IMultiSendMessage
	GetInstrumentList() ([]InstrumentStatus, error)
	SubscribeFullMarketData(item string)
	UnsubscribeFullMarketData(item string)
}

type IFmdManagerService interface {
	IFmdManager
	IFxService.IFxServices
}

type IFmdManagerData interface {
	IFmdManager
	IDataShutDown.IDataShutDown
}
