package fullMarketDataHelper

type IFullMarketDataHelper interface {
	InstrumentListChannelName() string
	InstrumentChannelName(instrument string) string
	AllInstrumentChannelName() string
	FullMarketDataServiceInbound() string
	FullMarketDataServicePublishInstrumentList() string
	RegisteredSource(instrument string) string
}
