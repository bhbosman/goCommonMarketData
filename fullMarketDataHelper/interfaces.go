package fullMarketDataHelper

import "github.com/bhbosman/gocommon/messageRouter"

type IFullMarketDataHelper interface {
	InstrumentListChannelName() string
	InstrumentChannelName(instrument string) string
	AllInstrumentChannelName() string
	FullMarketDataServiceInbound() string
	FullMarketDataServicePublishInstrumentList() string
	RegisteredSource(instrument string) string
	RegisterConsumerMethods(router messageRouter.MessageRouter)
}
