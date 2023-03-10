package fullMarketDataManagerViewer

import (
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	ui2 "github.com/bhbosman/goUi/ui"
	"github.com/rivo/tview"
)

type factory struct {
	Service                    IFullMarketDataViewService
	app                        *tview.Application
	instrumentReferenceService instrumentReference.IInstrumentReferenceService
}

func (self *factory) OrderNumber() int {
	return 2
}

func (self *factory) Content() (string, ui2.IPrimitiveCloser, error) {
	slide, err := newSlide(
		self.app,
		self.instrumentReferenceService,
		self.Service,
	)
	if err != nil {
		return "", nil, err
	}
	return self.Title(), slide, nil
}

func (self *factory) Title() string {
	return "MarketData"
}

func NewCoverSlideFactory(
	Service IFullMarketDataViewService,
	app *tview.Application,
	instrumentReferenceService instrumentReference.IInstrumentReferenceService,
) *factory {
	return &factory{
		Service:                    Service,
		app:                        app,
		instrumentReferenceService: instrumentReferenceService,
	}
}
