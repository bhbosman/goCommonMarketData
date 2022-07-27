package instrumentReference

import (
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
)

type MarketDataFeedReference struct {
	LunoFeeds   []*LunoReferenceData
	KrakenFeeds []*KrakenReferenceData
}

type IReferenceData interface {
	PriceDecimals() int
	VolumeDecimals() int
}

type ReferenceData struct {
	priceDecimals  int
	volumeDecimals int
}

func (self *ReferenceData) PriceDecimals() int {
	return self.priceDecimals
}

func (self *ReferenceData) VolumeDecimals() int {
	return self.volumeDecimals
}

type LunoReferenceData struct {
	ReferenceData
	SystemName       string
	Provider         string
	Name             string
	MappedInstrument string
}

type KrakenFeed struct {
	ReferenceData
	SystemName       string
	Pair             string
	Type             string
	MappedInstrument string
}

type KrakenReferenceData struct {
	ConnectionName string
	Provider       string
	Feeds          []*KrakenFeed
}

func NewDefaultReferenceData() *LunoReferenceData {
	return &LunoReferenceData{
		ReferenceData: ReferenceData{
			priceDecimals:  6,
			volumeDecimals: 6,
		},
		SystemName:       "Default.Default",
		Provider:         "Default",
		Name:             "Default",
		MappedInstrument: "",
	}
}

type IInstrumentReference interface {
	ISendMessage.ISendMessage
	SomeMethod()
	GetReferenceData(instrumentData string) (IReferenceData, bool)
	GetLunoProviders() ([]LunoReferenceData, error)
	GetKrakenProviders() ([]KrakenReferenceData, error)
}

type IInstrumentReferenceService interface {
	IInstrumentReference
	IFxService.IFxServices
}

type IInstrumentReferenceData interface {
	IInstrumentReference
	IDataShutDown.IDataShutDown
	ISendMessage.ISendMessage
}
