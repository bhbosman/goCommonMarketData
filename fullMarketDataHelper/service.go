package fullMarketDataHelper

import (
	"fmt"
	"github.com/bhbosman/gocommon/messageRouter"
)

const FullMarketDataServiceInbound = "FullMarketDataService_Inbound"
const FullMarketDataServicePublishAll = "FM_ALL"
const FullMarketDataServicePublishInstrumentList = "FMD_InstrumentList"

type FullMarketDataHelper struct {
}

func (self *FullMarketDataHelper) RegisterConsumerMethods(router messageRouter.MessageRouter) {
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

func NewFullMarketDataHelper() (IFullMarketDataHelper, error) {
	return &FullMarketDataHelper{}, nil
}
