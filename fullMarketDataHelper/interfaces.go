package fullMarketDataHelper

type IFullMarketDataHelper interface {
	InstrumentListChannelName() string
	InstrumentChannelName(instrument string) string
	InstrumentChannelNameMulti(instruments ...string) []string
	InstrumentChannelNameForTop5(instrument string) string
	InstrumentChannelNameForTop5Multi(instruments ...string) []string
	FullMarketDataServicePublishInstrumentList() string
	RegisteredSource(instrument string) string
}
