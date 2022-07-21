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
	OrderSide          [2]*avltree.Tree
	SourceTimestamp    int64
	SourceMessageCount int64
	messageRouter      *messageRouter.MessageRouter
	Name               string
}

func (self *FullMarketOrderBook) clear() {
	self.Orders = make(map[string]struct {
		stream.OrderSide
		*PricePoint
	})
	self.SourceTimestamp = 0
	self.OrderSide[0].Clear()
	self.OrderSide[1].Clear()
}

func (self *FullMarketOrderBook) addOrder(order *FullMarketOrder) {
	get, found := self.OrderSide[order.side].Get(order.price)
	if found {
		if pricePoint, ok := get.(*PricePoint); ok {
			pricePoint.AddOrder(order)
			self.Orders[order.id] = struct {
				stream.OrderSide
				*PricePoint
			}{
				order.side,
				pricePoint,
			}
		}
	} else {
		pricePoint := NewPricePoint(order.price)
		pricePoint.AddOrder(order)
		self.OrderSide[order.side].Put(order.price, pricePoint)
		self.Orders[order.id] = struct {
			stream.OrderSide
			*PricePoint
		}{
			order.side,
			pricePoint,
		}
	}
}

var epsilon = 1e-8

func (self *FullMarketOrderBook) tradeUpdate(tradeUpdate *stream.FullMarketData_ReduceVolumeInstruction) {
	if makerOrder, ok := self.Orders[tradeUpdate.Id]; ok {
		if find, order := makerOrder.PricePoint.Find(tradeUpdate.Id); find {
			newVolume := order.ReduceVolume(tradeUpdate.Volume)
			if newVolume <= epsilon {
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

func (self *FullMarketOrderBook) SetTimeStamp(timestamp int64) {
	self.SourceTimestamp = timestamp
}

func (self *FullMarketOrderBook) UpdateMessageReceivedCount() {
	self.SourceMessageCount++
}

func (self *FullMarketOrderBook) handleFullMarketDataAddOrder(msg *stream.FullMarketData_AddOrderInstruction) {
	order := newFullMarketOrder(msg.Order.Side, msg.Order.Id, msg.Order.Price, msg.Order.Volume)
	self.addOrder(order)

}

func (self *FullMarketOrderBook) handleFullMarketDataClear(*stream.FullMarketData_Clear) {
	self.clear()
}

func (self *FullMarketOrderBook) handleFullMarketDataReduceVolume(msg *stream.FullMarketData_ReduceVolumeInstruction) {
	self.tradeUpdate(msg)
}

func (self *FullMarketOrderBook) handleFullMarketDataRemoveInstrumentInstruction(order *stream.FullMarketData_RemoveInstrumentInstruction) {

}

func (self *FullMarketOrderBook) handleFullMarketDataDeleteOrder(order *stream.FullMarketData_DeleteOrderInstruction) {
	self.deleteOrder(order.Id)
}

func (self *FullMarketOrderBook) Send(msg interface{}) {
	self.messageRouter.Route(msg)
}
func NewFullMarketOrderBook(name string) *FullMarketOrderBook {
	result := &FullMarketOrderBook{
		Name: name,
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
	result.messageRouter.Add(result.handleFullMarketDataClear)
	result.messageRouter.Add(result.handleFullMarketDataAddOrder)
	result.messageRouter.Add(result.handleFullMarketDataReduceVolume)
	result.messageRouter.Add(result.handleFullMarketDataDeleteOrder)
	result.messageRouter.Add(result.handleFullMarketDataRemoveInstrumentInstruction)
	return result
}
