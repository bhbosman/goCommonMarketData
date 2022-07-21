package fullMarketDataManagerService

import (
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
)

type IFmdManager interface {
	GetInstrumentList() ([]string, error)
}

type IFmdManagerService interface {
	IFmdManager
	IFxService.IFxServices
}

type IFmdManagerData interface {
	IFmdManager
	ISendMessage.ISendMessage
	IDataShutDown.IDataShutDown
}
