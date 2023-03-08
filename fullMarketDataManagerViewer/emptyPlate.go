package fullMarketDataManagerViewer

import (
	"github.com/rivo/tview"
)

type emptyCell struct {
}

func (self *emptyCell) GetCell(row, column int) *tview.TableCell {
	return tview.NewTableCell("")
}

func (self *emptyCell) GetRowCount() int {
	return 0
}

func (self *emptyCell) GetColumnCount() int {
	return 0
}

func (self *emptyCell) SetCell(row, column int, cell *tview.TableCell) {
}

func (self *emptyCell) RemoveRow(row int) {
}

func (self *emptyCell) RemoveColumn(column int) {
}

func (self *emptyCell) InsertRow(row int) {
}

func (self *emptyCell) InsertColumn(column int) {
}

func (self *emptyCell) Clear() {
}
