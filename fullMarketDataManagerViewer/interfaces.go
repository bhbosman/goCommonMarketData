package fullMarketDataManagerViewer

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
)

type IFullMarketDataView interface {
	ISendMessage.ISendMessage
	SetMarketDataListChange(change func(list []fullMarketDataManagerService.InstrumentStatus) bool)
	SetMarketDataInstanceChange(change func(data *stream.PublishTop5) bool)
	UnsubscribeFullMarketData(item string)
	SubscribeFullMarketData(item string)
}

type IFullMarketDataViewData interface {
	IFullMarketDataView
	IDataShutDown.IDataShutDown
	Start(serviceName string)
}
type IFullMarketDataViewService interface {
	IFullMarketDataView
	IFxService.IFxServices
}
