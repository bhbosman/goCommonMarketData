package fullMarketDataManagerViewer

import (
	"github.com/rivo/tview"
	"strconv"
)

type listPlate struct {
	list      []string
	emptyCell *tview.TableCell
}

func (self *listPlate) GetCell(row, column int) *tview.TableCell {
	if row == -1 || column == -1 {
		return tview.NewTableCell("")
	}

	switch row {
	case 0:
		switch column {
		case 0:
			return tview.NewTableCell("*").SetSelectable(false).SetAlign(tview.AlignRight)
		case 1:
			return tview.NewTableCell("Name").SetSelectable(false)
		}
	default:
		switch column {
		case 0:
			return tview.NewTableCell(strconv.Itoa(row)).SetSelectable(false).SetAlign(tview.AlignRight)
		case 1:
			n := row - 1
			c := len(self.list)
			if c > n {
				return tview.NewTableCell(self.list[row-1])
			}
		}
	}

	return self.emptyCell
}

func (self *listPlate) GetRowCount() int {
	return len(self.list) + 1
}

func (self *listPlate) GetColumnCount() int {
	return 2
}

func (self *listPlate) SetCell(row, column int, cell *tview.TableCell) {
}

func (self *listPlate) RemoveRow(row int) {
}

func (self *listPlate) RemoveColumn(column int) {
}

func (self *listPlate) InsertRow(row int) {
}

func (self *listPlate) InsertColumn(column int) {
}

func (self *listPlate) Clear() {
}

func (self *listPlate) GetItem(row int) (string, bool) {

	if row == -1 {
		return "", false
	}

	index := row - 1
	count := len(self.list)
	if index >= 0 && count > index {
		return self.list[index], true
	}
	return "", false
}

func newListPlate(list []string) *listPlate {
	return &listPlate{
		list:      list,
		emptyCell: tview.NewTableCell("").SetSelectable(false),
	}
}
