package fullMarketDataManagerViewer

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/goUi/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type slide struct {
	slideOrderNumber           int
	service                    IFullMarketDataViewService
	next                       tview.Primitive
	canDraw                    bool
	app                        *tview.Application
	table                      *tview.Table
	listTable                  *tview.Table
	marketDataListPlate        *listPlate
	marketDataPlate            *marketDataPlate
	instrumentReferenceService instrumentReference.IInstrumentReferenceService
	currentReferenceData       instrumentReference.IReferenceData
	selectedItem               string
	slideName                  string
}

func (self *slide) OrderNumber() int {
	return self.slideOrderNumber
}

func (self *slide) Name() string {
	return self.slideName
}

func (self *slide) Toggle(b bool) {
	self.canDraw = b
	switch b {
	case true:
		if self.selectedItem != "" {
			self.service.SubscribeFullMarketData(self.selectedItem)
		}
	case false:
		if self.selectedItem != "" {
			self.service.UnsubscribeFullMarketData(self.selectedItem)
		}
	}
	if b {
		self.app.ForceDraw()
	}
}

func (self *slide) Draw(screen tcell.Screen) {
	self.next.Draw(screen)
}

func (self *slide) GetRect() (int, int, int, int) {
	return self.next.GetRect()
}

func (self *slide) SetRect(x, y, width, height int) {
	self.next.SetRect(x, y, width, height)
}

func (self *slide) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return self.next.InputHandler()
}

func (self *slide) Focus(delegate func(p tview.Primitive)) {
	self.next.Focus(delegate)
}

func (self *slide) HasFocus() bool {
	return self.next.HasFocus()
}

func (self *slide) Blur() {
	self.next.Blur()
}

func (self *slide) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return self.next.MouseHandler()
}

func (self *slide) Close() error {
	return nil
}

func (self *slide) UpdateContent() error {
	return nil
}

func (self *slide) init() {
	self.table = tview.NewTable()
	self.table.SetSelectable(true, false)
	self.table.SetBorder(true)
	self.table.SetFixed(1, 1)
	self.table.SetTitle("Full Market Data Viewer")
	self.table.SetContent(&emptyCell{})

	self.listTable = tview.NewTable()
	self.listTable.SetBorder(true)
	self.listTable.SetSelectable(true, false)
	self.listTable.SetFixed(1, 1)
	self.listTable.SetSelectionChangedFunc(
		func(row, column int) {
			row, _ = self.listTable.GetSelection()
			if item, ok := self.marketDataListPlate.GetItem(row); ok {
				if item != self.selectedItem {
					if self.selectedItem != "" {
						self.service.UnsubscribeFullMarketData(self.selectedItem)
					}
					self.selectedItem = item
					if self.canDraw {
						self.service.SubscribeFullMarketData(self.selectedItem)
					}
					self.currentReferenceData, _ = self.instrumentReferenceService.GetReferenceData(item)
					if self.currentReferenceData == nil {
						self.currentReferenceData = instrumentReference.NewDefaultReferenceData()
					}
				}
			}
		},
	)
	self.listTable.SetSelectedFunc(
		func(row, column int) {

		},
	)
	self.listTable.SetContent(&emptyCell{})
	flex := tview.NewFlex().
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(self.listTable, 30, 1, true).
				AddItem(self.table, 0, 3, false),
			0,
			1,
			true)

	self.next = flex
}

func (self *slide) OnSetMarketDataListChange(list []fullMarketDataManagerService.InstrumentStatus) bool {
	return self.app.QueueUpdate(
		func() {
			if list != nil && len(list) > 0 {
				plateNil := self.marketDataListPlate == nil
				self.marketDataListPlate = newListPlate(list)
				self.listTable.SetContent(self.marketDataListPlate)
				if plateNil && self.listTable != nil && len(self.marketDataListPlate.list) > 0 {
					self.listTable.Select(1, 0)
				} else {
					row, column := self.listTable.GetSelection()
					self.listTable.Select(row, column)
				}
			} else {
				if self.selectedItem != "" {
					self.service.UnsubscribeFullMarketData(self.selectedItem)
				}
				self.selectedItem = ""
				self.marketDataListPlate = nil
				self.listTable.SetContent(&emptyCell{})
				self.marketDataPlate = nil
				self.table.SetContent(&emptyCell{})
			}
			if self.canDraw {
				self.app.ForceDraw()
			}
		},
	)
}

func (self *slide) OnSetMarketDataInstanceChange(data *stream.PublishTop5) bool {
	return self.app.QueueUpdate(
		func() {
			row, _ := self.listTable.GetSelection()
			if text, ok := self.marketDataListPlate.GetItem(row); ok {
				if text == data.Instrument {
					if self.currentReferenceData == nil {
						self.currentReferenceData = instrumentReference.NewDefaultReferenceData()
					}
					self.marketDataPlate = newMarketDataPlate(data, self.currentReferenceData)
					self.table.SetContent(self.marketDataPlate)
					self.table.ScrollToBeginning()
					if self.canDraw {
						self.app.ForceDraw()
					}
				}
			}
		},
	)
}

func newSlide(
	slideOrderNumber int,
	slideName string,
	app *tview.Application,
	instrumentReferenceService instrumentReference.IInstrumentReferenceService,
	service IFullMarketDataViewService,
) (ui.IPrimitiveCloser, error) {
	result := &slide{
		slideOrderNumber:           slideOrderNumber,
		slideName:                  slideName,
		app:                        app,
		service:                    service,
		instrumentReferenceService: instrumentReferenceService,
	}
	result.init()
	service.SetMarketDataListChange(result.OnSetMarketDataListChange)
	service.SetMarketDataInstanceChange(result.OnSetMarketDataInstanceChange)
	return result, nil
}
