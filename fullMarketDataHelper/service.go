package fullMarketDataHelper

import (
	"fmt"
	"github.com/cskr/pubsub"
)

const FullMarketDataServicePublishInstrumentList = "FMD_InstrumentList"

type FullMarketDataHelper struct {
	pubSub *pubsub.PubSub
}

func (self *FullMarketDataHelper) InstrumentChannelNameForTop5Multi(instruments ...string) []string {
	result := make([]string, len(instruments))
	for i, s := range instruments {
		result[i] = self.InstrumentChannelNameForTop5(s)
	}
	return result
}

func (self *FullMarketDataHelper) InstrumentChannelNameMulti(instruments ...string) []string {
	result := make([]string, len(instruments))
	for i, s := range instruments {
		result[i] = self.InstrumentChannelName(s)
	}
	return result
}

func (self *FullMarketDataHelper) RegisteredSource(instrument string) string {
	return fmt.Sprintf("FMD.RegisteredSource.INSTANCE.%v", instrument)
}

func (self *FullMarketDataHelper) FullMarketDataServicePublishInstrumentList() string {
	return FullMarketDataServicePublishInstrumentList
}

func (self *FullMarketDataHelper) InstrumentListChannelName() string {
	return self.instrumentListChannelName()
}

func (self *FullMarketDataHelper) InstrumentChannelName(instrument string) string {
	return self.instrumentChannelName(instrument)
}

func (self *FullMarketDataHelper) instrumentListChannelName() string {
	return FullMarketDataServicePublishInstrumentList
}

func (self *FullMarketDataHelper) instrumentChannelName(instrument string) string {
	return fmt.Sprintf("FMD.INSTANCE.%v", instrument)
}
func (self *FullMarketDataHelper) InstrumentChannelNameForTop5(instrument string) string {
	return fmt.Sprintf("FMD.CALCULATED.FULL.ORDER.BOOK.INSTANCE.%v", instrument)
}

func NewFullMarketDataHelper(pubSub *pubsub.PubSub,
) (IFullMarketDataHelper, error) {
	return &FullMarketDataHelper{
		pubSub: pubSub,
	}, nil
}
