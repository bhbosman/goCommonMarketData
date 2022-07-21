package fullMarketData

import (
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"github.com/emirpasic/gods/trees/avltree"
)

type IFullMarketOrderBook interface {
	ISendMessage.ISendMessage
	OrderCount() int
	BidOrderSide() *avltree.Tree
	AskOrderSide() *avltree.Tree
	InstrumentName() string
}
