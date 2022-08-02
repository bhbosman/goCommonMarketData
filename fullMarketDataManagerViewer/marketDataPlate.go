package fullMarketDataManagerViewer

import (
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goMessages/marketData/stream"
	"github.com/rivo/tview"
	"strconv"
)

type marketDataPlate struct {
	data                 *stream.PublishTop5
	emptyCell            *tview.TableCell
	currentReferenceData instrumentReference.IReferenceData
}

func newMarketDataPlate(
	data *stream.PublishTop5,
	currentReferenceData instrumentReference.IReferenceData,
) *marketDataPlate {
	return &marketDataPlate{
		data:                 data,
		emptyCell:            tview.NewTableCell("").SetSelectable(false),
		currentReferenceData: currentReferenceData,
	}
}

func (self *marketDataPlate) GetCell(row, column int) *tview.TableCell {
	if row == -1 || column == -1 {
		return tview.NewTableCell("")
	}

	switch row {
	case 0:
		switch column {
		case 0:
			return tview.NewTableCell("*").SetSelectable(false)
		case 1:
			return tview.NewTableCell("Count").
				SetSelectable(false).
				SetMaxWidth(20).
				SetAlign(tview.AlignRight)
		case 2:
			return tview.NewTableCell("BidVol").
				SetSelectable(false).
				SetMaxWidth(20).
				SetAlign(tview.AlignRight)
		case 3:
			return tview.NewTableCell("Bid").
				SetSelectable(false).
				SetMaxWidth(20).
				SetAlign(tview.AlignRight)
		case 4:
			return tview.NewTableCell("Ask").
				SetSelectable(false).
				SetMaxWidth(20).
				SetAlign(tview.AlignRight)
		case 5:
			return tview.NewTableCell("AskVol").
				SetSelectable(false).
				SetMaxWidth(20).
				SetAlign(tview.AlignRight)
		case 6:
			return tview.NewTableCell("Count").
				SetSelectable(false).
				SetMaxWidth(20).
				SetAlign(tview.AlignRight)
		case 7:
			return tview.NewTableCell("Split").
				SetSelectable(false).
				SetMaxWidth(20).
				SetAlign(tview.AlignRight)
		}
	default:
		switch column {
		case 0:
			return tview.NewTableCell(strconv.Itoa(row)).SetSelectable(false).SetAlign(tview.AlignRight)
		case 1:
			index := row - 1
			count := len(self.data.Bid)
			if index >= 0 && index < count {
				return tview.NewTableCell(strconv.Itoa(int(self.data.Bid[index].OpenOrderCount))).SetSelectable(false).SetAlign(tview.AlignRight)
			}
			break
		case 2:
			index := row - 1
			count := len(self.data.Bid)
			if index >= 0 && index < count {
				s := strconv.FormatFloat(self.data.Bid[index].Volume, 'f', self.currentReferenceData.VolumeDecimals(), 32)
				return tview.NewTableCell(s).
					SetSelectable(true).
					SetMaxWidth(20).
					SetAlign(tview.AlignRight)

			}
			break
		case 3:
			index := row - 1
			count := len(self.data.Bid)
			if index >= 0 && index < count {
				s := strconv.FormatFloat(self.data.Bid[index].Price, 'f', self.currentReferenceData.PriceDecimals(), 32)
				return tview.NewTableCell(s).
					SetSelectable(true).
					SetMaxWidth(20).
					SetAlign(tview.AlignRight)
			}
			break
		case 4:
			index := row - 1
			count := len(self.data.Ask)
			if index >= 0 && index < count {
				s := strconv.FormatFloat(self.data.Ask[index].Price, 'f', self.currentReferenceData.PriceDecimals(), 32)
				return tview.NewTableCell(s).
					SetSelectable(true).
					SetMaxWidth(20).
					SetAlign(tview.AlignRight)
			}
			break
		case 5:
			index := row - 1
			count := len(self.data.Ask)
			if index >= 0 && index < count {
				s := strconv.FormatFloat(self.data.Ask[index].Volume, 'f', self.currentReferenceData.VolumeDecimals(), 32)
				return tview.NewTableCell(s).
					SetSelectable(true).
					SetMaxWidth(20).
					SetAlign(tview.AlignRight)
			}
			break
		case 6:
			index := row - 1
			count := len(self.data.Ask)
			if index >= 0 && index < count {
				return tview.NewTableCell(strconv.Itoa(int(self.data.Ask[index].OpenOrderCount))).
					SetSelectable(false).
					SetAlign(tview.AlignRight)
			}
			break
		case 7:
			if row == 1 {
				if len(self.data.Bid) > 0 && len(self.data.Ask) > 0 {
					v := self.data.Ask[0].Price - self.data.Bid[0].Price
					s := strconv.FormatFloat(v, 'f', self.currentReferenceData.PriceDecimals(), 32)
					return tview.NewTableCell(s).
						SetSelectable(true).
						SetMaxWidth(20).
						SetAlign(tview.AlignRight)
				}
			}
			break
		}
	}
	return self.emptyCell
}

func (self *marketDataPlate) dataCount() int {
	if len(self.data.Bid) >= len(self.data.Ask) {
		return len(self.data.Bid)
	} else {
		return len(self.data.Ask)
	}

}
func (self *marketDataPlate) GetRowCount() int {
	return self.dataCount() + 1
}

func (self *marketDataPlate) GetColumnCount() int {
	return 8
}

func (self *marketDataPlate) SetCell(row, column int, cell *tview.TableCell) {
}

func (self *marketDataPlate) RemoveRow(row int) {
}

func (self *marketDataPlate) RemoveColumn(column int) {
}

func (self *marketDataPlate) InsertRow(row int) {
}

func (self *marketDataPlate) InsertColumn(column int) {
}

func (self *marketDataPlate) Clear() {
}
