package fullMarketData

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketData/stream"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/emirpasic/gods/utils"
)

const BuySide = stream.OrderSide_BidOrder
const AskSide = stream.OrderSide_AskOrder

type FullMarketOrderBook struct {
	Orders map[string]struct {
		stream.OrderSide
		*PricePoint
	}
	OrderSide     [2]*avltree.Tree
	messageRouter messageRouter.IMessageRouter
	feedName      string
	name          string
	status        string
}

func (self *FullMarketOrderBook) Status() string {
	return self.status
}

func (self *FullMarketOrderBook) SetStatus(status string) {
	self.status = status
}

func (self *FullMarketOrderBook) FeedName() string {
	return self.feedName
}

func (self *FullMarketOrderBook) InstrumentName() string {
	return self.name
}

func (self *FullMarketOrderBook) BidOrderSide() *avltree.Tree {
	return self.OrderSide[stream.OrderSide_BidOrder]
}

func (self *FullMarketOrderBook) AskOrderSide() *avltree.Tree {
	return self.OrderSide[stream.OrderSide_AskOrder]
}

func (self *FullMarketOrderBook) OrderCount() int {
	return len(self.Orders)
}

func (self *FullMarketOrderBook) clear() {
	self.Orders = make(map[string]struct {
		stream.OrderSide
		*PricePoint
	})
	self.OrderSide[0].Clear()
	self.OrderSide[1].Clear()
}

func (self *FullMarketOrderBook) addOrder(order *FullMarketOrder) {
	get, found := self.OrderSide[order.Side].Get(order.Price)
	if found {
		if pricePoint, ok := get.(*PricePoint); ok {
			pricePoint.AddOrder(order)
			self.Orders[order.Id] = struct {
				stream.OrderSide
				*PricePoint
			}{
				order.Side,
				pricePoint,
			}
		}
	} else {
		pricePoint := NewPricePoint(order.Price)
		pricePoint.AddOrder(order)
		self.OrderSide[order.Side].Put(order.Price, pricePoint)
		self.Orders[order.Id] = struct {
			stream.OrderSide
			*PricePoint
		}{
			order.Side,
			pricePoint,
		}
	}
}

var epsilon = 1e-9

func (self *FullMarketOrderBook) tradeUpdate(tradeUpdate *stream.FullMarketData_ReduceVolumeInstruction) {
	if makerOrder, ok := self.Orders[tradeUpdate.Id]; ok {
		if find, order := makerOrder.PricePoint.Find(tradeUpdate.Id); find {
			d := tradeUpdate.Volume
			newVolume := order.ReduceVolume(d)
			if newVolume < epsilon {
				self.deleteOrder(tradeUpdate.Id)
			}
		}
	}
}

func (self *FullMarketOrderBook) deleteOrder(oderId string) {
	if data, ok := self.Orders[oderId]; ok {
		delete(self.Orders, oderId)
		data.PricePoint.Delete(oderId)
		if data.PricePoint.Count() == 0 {
			self.OrderSide[data.OrderSide].Remove(data.PricePoint.Price)
		}
	}
}

func (self *FullMarketOrderBook) handleFullMarketDataAddOrder(msg *stream.FullMarketData_AddOrderInstruction) {
	order := newFullMarketOrder(msg.Order.Side, msg.Order.Id, msg.Order.Price, msg.Order.Volume, msg.Order.ExtraData)
	self.addOrder(order)

}

func (self *FullMarketOrderBook) handleFullMarketDataClear(*stream.FullMarketData_Clear) {
	self.clear()
}

func (self *FullMarketOrderBook) handleFullMarketDataReduceVolume(msg *stream.FullMarketData_ReduceVolumeInstruction) {
	self.tradeUpdate(msg)
}

func (self *FullMarketOrderBook) handleFullMarketData_Instrument_InstrumentStatus(msg *stream.FullMarketData_Instrument_InstrumentStatus) {
	self.status = msg.Status
}

func (self *FullMarketOrderBook) handleFullMarketDataRemoveInstrumentInstruction(_ *stream.FullMarketData_RemoveInstrumentInstruction) {

}

func (self *FullMarketOrderBook) handleFullMarketDataDeleteOrder(order *stream.FullMarketData_DeleteOrderInstruction) {
	self.deleteOrder(order.Id)
}

func (self *FullMarketOrderBook) Send(msg interface{}) error {
	self.messageRouter.Route(msg)
	return nil
}

func NewFullMarketOrderBook(feedName, name string) IFullMarketOrderBook {
	result := &FullMarketOrderBook{
		name:     name,
		feedName: feedName,
		Orders: make(map[string]struct {
			stream.OrderSide
			*PricePoint
		}),
		OrderSide: [2]*avltree.Tree{
			avltree.NewWith(utils.Float64Comparator),
			avltree.NewWith(utils.Float64Comparator),
		},
		messageRouter: messageRouter.NewMessageRouter(),
	}
	_ = result.messageRouter.Add(result.handleFullMarketDataClear)
	_ = result.messageRouter.Add(result.handleFullMarketDataAddOrder)
	_ = result.messageRouter.Add(result.handleFullMarketDataReduceVolume)
	_ = result.messageRouter.Add(result.handleFullMarketDataDeleteOrder)
	_ = result.messageRouter.Add(result.handleFullMarketData_Instrument_InstrumentStatus)
	_ = result.messageRouter.Add(result.handleFullMarketDataRemoveInstrumentInstruction)

	return result
}
