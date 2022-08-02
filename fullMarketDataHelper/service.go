package fullMarketDataHelper

import (
	"fmt"
	"github.com/cskr/pubsub"
)

const FullMarketDataServiceInbound = "FullMarketDataService_Inbound"
const FullMarketDataServicePublishAll = "FM_ALL"
const FullMarketDataServicePublishInstrumentList = "FMD_InstrumentList"

type FullMarketDataHelper struct {
	pubSub *pubsub.PubSub
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

func (self *FullMarketDataHelper) FullMarketDataServiceInbound() string {
	return FullMarketDataServiceInbound
}

func (self *FullMarketDataHelper) InstrumentListChannelName() string {
	return self.instrumentListChannelName()
}

func (self *FullMarketDataHelper) InstrumentChannelName(instrument string) string {
	return self.instrumentChannelName(instrument)
}

func (self *FullMarketDataHelper) AllInstrumentChannelName() string {
	return self.allInstrumentChannelName()
}
func (self *FullMarketDataHelper) instrumentListChannelName() string {
	return FullMarketDataServicePublishInstrumentList
}

func (self *FullMarketDataHelper) allInstrumentChannelName() string {
	return FullMarketDataServicePublishAll
}

func (self *FullMarketDataHelper) instrumentChannelName(instrument string) string {
	return fmt.Sprintf("FMD.INSTANCE.%v", instrument)
}

func NewFullMarketDataHelper(pubSub *pubsub.PubSub,
) (IFullMarketDataHelper, error) {
	return &FullMarketDataHelper{
		pubSub: pubSub,
	}, nil
}
