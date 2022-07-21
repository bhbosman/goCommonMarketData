package fullMarketDataManagerService

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketData"
	stream2 "github.com/bhbosman/goCommonMarketData/fullMarketData/stream"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/cskr/pubsub"
	"sort"
)

type data struct {
	outstandingRequests      []interface{}
	queuedMessages           map[string]*stream.PublishTop5
	MessageRouter            *messageRouter.MessageRouter
	pubSub                   *pubsub.PubSub
	fmd                      map[string]*fullMarketData.FullMarketOrderBook
	fmdDirty                 map[string]bool
	publishListOfInstruments bool
	fullMarketDataHelper     fullMarketDataHelper.IFullMarketDataHelper
}

func (self *data) GetInstrumentList() ([]string, error) {
	return self.buildInstrumentList()
}

func (self *data) Send(message interface{}) error {
	_, err := self.MessageRouter.Route(message)
	return err
}

func (self *data) SomeMethod() {
}

func (self *data) ShutDown() error {
	return nil
}

func (self *data) findFullMarketDataBook(name string) *fullMarketData.FullMarketOrderBook {
	if result, ok := self.fmd[name]; ok {
		return result
	}
	result := fullMarketData.NewFullMarketOrderBook(name)
	self.fmd[name] = result
	self.publishListOfInstruments = true
	return result
}

func (self *data) handleFullMarketDataInstrumentListRequest(request *stream2.FullMarketData_InstrumentListRequest) {

}

func (self *data) handleFullMarketDataRemoveInstrumentInstruction(msg *stream2.FullMarketData_RemoveInstrumentInstruction) {
	delete(self.fmd, msg.Instrument)
	self.publishListOfInstruments = true
}

func (self *data) handlePublishFullMarketData(msg *PublishFullMarketData) {
	self.outstandingRequests = append(self.outstandingRequests, msg)
}

func (self *data) handleEmptyQueue(msg *messages.EmptyQueue) {
	// Publish List
	if self.publishListOfInstruments {
		ss, _ := self.buildInstrumentList()
		if self.pubSub.PubWithContext(
			&stream2.FullMarketData_InstrumentListResponse{
				Instruments: ss,
			},
			self.fullMarketDataHelper.InstrumentListChannelName(),
		) {
			self.publishListOfInstruments = false
		} else {
			msg.ErrorHappen = true
			return
		}
	}

	for _, request := range self.outstandingRequests {
		switch req := request.(type) {
		case *PublishFullMarketData:
			if v, ok := self.fmd[req.PublishInstrument]; ok {
				top5, b := self.calculate(true, v)
				if b {
					self.queuedMessages[req.PublishInstrument] = top5
				}
			}
		}
	}
	self.outstandingRequests = nil

	for key := range self.fmdDirty {
		if v, ok := self.fmd[key]; ok {
			top5, b := self.calculate(false, v)
			if b {
				self.queuedMessages[key] = top5
			}
		}
	}

	self.fmdDirty = make(map[string]bool)

	var published []string
	publishError := false
	for key, value := range self.queuedMessages {

		ss := []string{
			self.fullMarketDataHelper.AllInstrumentChannelName(),
			self.fullMarketDataHelper.InstrumentChannelName(key),
		}
		if self.pubSub.PubWithContext(value, ss...) {
			published = append(published, key)
		} else {
			publishError = true
			msg.ErrorHappen = true
			return
		}
	}
	if !publishError {
		self.queuedMessages = make(map[string]*stream.PublishTop5)
	} else {
		for _, s := range published {
			delete(self.queuedMessages, s)
		}
	}
}

func (self *data) handleFullMarketDataAddOrder(msg *stream2.FullMarketData_AddOrderInstruction) {
	self.findFullMarketDataBook(msg.Instrument).Send(msg)
	self.fmdDirty[msg.Instrument] = true
}

func (self *data) handleFullMarketDataClear(msg *stream2.FullMarketData_Clear) {
	self.findFullMarketDataBook(msg.Instrument).Send(msg)
	self.fmdDirty[msg.Instrument] = true
}

func (self *data) handleFullMarketDataReduceVolume(msg *stream2.FullMarketData_ReduceVolumeInstruction) {
	self.findFullMarketDataBook(msg.Instrument).Send(msg)
	self.fmdDirty[msg.Instrument] = true
}

func (self *data) handleFullMarketDataDeleteOrder(msg *stream2.FullMarketData_DeleteOrderInstruction) {
	self.findFullMarketDataBook(msg.Instrument).Send(msg)
	self.fmdDirty[msg.Instrument] = true
}

func (self *data) calculate(force bool, fullMarketOrderBook *fullMarketData.FullMarketOrderBook) (*stream.PublishTop5, bool) {

	thereWasAChange := force || len(fullMarketOrderBook.Orders) == 0
	var bids []*stream.Point
	if highBidNode := fullMarketOrderBook.OrderSide[fullMarketData.BuySide].Right(); highBidNode != nil {
		count := 0
		for node := highBidNode; node != nil; node = node.Prev() {
			bidPrice := node.Key.(float64)
			if pp, ok := node.Value.(*fullMarketData.PricePoint); ok {
				thereWasAChange = thereWasAChange || pp.Touched
				pp.ClearTouched()
				volume := pp.GetVolume()
				bids = append(bids, &stream.Point{
					Price:          bidPrice,
					Volume:         volume,
					OpenOrderCount: int32(pp.Count()),
				})
			}
			count++
		}
	}
	var asks []*stream.Point
	if lowAskNode := fullMarketOrderBook.OrderSide[fullMarketData.AskSide].Left(); lowAskNode != nil {
		count := 0
		for node := lowAskNode; node != nil; node = node.Next() {
			askPrice := node.Key.(float64)
			if pp, ok := node.Value.(*fullMarketData.PricePoint); ok {
				thereWasAChange = thereWasAChange || pp.Touched
				pp.ClearTouched()
				volume := pp.GetVolume()
				asks = append(asks, &stream.Point{
					Price:          askPrice,
					Volume:         volume,
					OpenOrderCount: int32(pp.Count()),
				})
			}
			count++
		}
	}
	spread := 0.0
	if len(asks) > 0 && len(bids) > 0 {
		spread = asks[0].Price - bids[0].Price
	}
	if thereWasAChange {
		return &stream.PublishTop5{
			Instrument:         fullMarketOrderBook.Name,
			Spread:             spread,
			SourceTimeStamp:    fullMarketOrderBook.SourceTimestamp,
			SourceMessageCount: fullMarketOrderBook.SourceMessageCount,
			Bid:                bids,
			Ask:                asks,
		}, true
	}
	return nil, false
}

func (self *data) buildInstrumentList() ([]string, error) {
	var ss []string
	for key := range self.fmd {
		ss = append(ss, key)
	}
	sort.Strings(ss)
	return ss, nil
}

func newData(
	pubSub *pubsub.PubSub,
	fullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper,
) (IFmdManagerData, error) {
	result := &data{
		outstandingRequests:  nil,
		queuedMessages:       make(map[string]*stream.PublishTop5),
		MessageRouter:        messageRouter.NewMessageRouter(),
		pubSub:               pubSub,
		fmd:                  make(map[string]*fullMarketData.FullMarketOrderBook),
		fmdDirty:             make(map[string]bool),
		fullMarketDataHelper: fullMarketDataHelper,
	}
	result.MessageRouter.Add(result.handleEmptyQueue)
	result.MessageRouter.Add(result.handleFullMarketDataClear)
	result.MessageRouter.Add(result.handleFullMarketDataAddOrder)
	result.MessageRouter.Add(result.handleFullMarketDataReduceVolume)
	result.MessageRouter.Add(result.handleFullMarketDataDeleteOrder)
	result.MessageRouter.Add(result.handlePublishFullMarketData)
	result.MessageRouter.Add(result.handleFullMarketDataRemoveInstrumentInstruction)
	result.MessageRouter.Add(result.handleFullMarketDataInstrumentListRequest)

	return result, nil
}
