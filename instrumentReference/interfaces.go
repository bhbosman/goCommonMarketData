package instrumentReference

import (
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
)

type ReferenceData struct {
	Name           string
	PriceDecimals  int
	VolumeDecimals int
}

func NewDefaultReferenceData() *ReferenceData {
	return &ReferenceData{
		Name:           "Default",
		PriceDecimals:  6,
		VolumeDecimals: 6,
	}
}

type IInstrumentReference interface {
	ISendMessage.ISendMessage
	SomeMethod()
	GetReferenceData(instrumentData string) (*ReferenceData, bool)
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
