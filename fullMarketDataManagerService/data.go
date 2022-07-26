package fullMarketDataManagerService

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketData"
	stream2 "github.com/bhbosman/goCommonMarketData/fullMarketData/stream"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/cskr/pubsub"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/reactivex/rxgo/v2"
	"google.golang.org/protobuf/proto"
	"sort"
)

type outstandingRequests struct {
	adder goCommsDefinitions.IAdder
	s     string
}

type data struct {
	proxy                    bool
	outstandingRequestsMap   map[outstandingRequests]interface{}
	queuedMessages           map[string]*stream.PublishTop5
	MessageRouter            *messageRouter.MessageRouter
	pubSub                   *pubsub.PubSub
	fmd                      map[string]fullMarketData.IFullMarketOrderBook
	fmdDirty                 map[string]bool
	fmdCount                 map[string]int
	publishListOfInstruments bool
	fullMarketDataHelper     fullMarketDataHelper.IFullMarketDataHelper
}

func (self *data) SubscribeFullMarketDataMulti(item ...string) {
	for _, s := range item {
		self.SubscribeFullMarketData(s)
	}
}

func (self *data) UnsubscribeFullMarketDataMulti(item ...string) {
	for _, s := range item {
		self.UnsubscribeFullMarketData(s)
	}
}

func (self *data) MultiSend(messages ...interface{}) {
	self.MessageRouter.MultiRoute(messages...)
}

func (self *data) SubscribeFullMarketData(item string) {
	if v, ok := self.fmdCount[item]; ok {
		self.fmdCount[item] = v + 1
	} else {
		self.fmdCount[item] = 1

		key := self.fullMarketDataHelper.RegisteredSource(item)
		self.pubSub.Pub(
			&stream2.FullMarketData_Instrument_Register{
				Instrument: item,
			},
			key,
		)
	}
}

func (self *data) UnsubscribeFullMarketData(item string) {
	if v, ok := self.fmdCount[item]; ok {
		v--
		self.fmdCount[item] = v
		if v == 0 {
			delete(self.fmdCount, item)
			if self.proxy {
				if book, ok := self.findFullMarketDataBook("", item); ok {
					clearMessage := &stream2.FullMarketData_Clear{
						Instrument: item,
					}
					_ = book.Send(clearMessage)
				}
			}
			key := self.fullMarketDataHelper.RegisteredSource(item)
			self.pubSub.Pub(
				&stream2.FullMarketData_Instrument_Unregister{
					Instrument: item,
				},
				key,
			)
		}
	}
}

func (self *data) GetInstrumentList() ([]InstrumentStatus, error) {
	return self.buildInstrumentList()
}

func (self *data) Send(message interface{}) error {
	self.MessageRouter.Route(message)
	return nil
}

func (self *data) SomeMethod() {
}

func (self *data) ShutDown() error {
	return nil
}

func (self *data) findFullMarketDataBook(feedName, name string) (fullMarketData.IFullMarketOrderBook, bool) {
	if result, ok := self.fmd[name]; ok {
		return result, true
	}
	result := fullMarketData.NewFullMarketOrderBook(feedName, name)
	self.fmd[name] = result
	self.publishListOfInstruments = true
	return result, true
}

func (self *data) handleFullMarketDataRemoveInstrumentInstruction(msg *stream2.FullMarketData_RemoveInstrumentInstruction) {
	if _, ok := self.fmd[msg.Instrument]; ok {
		self.MessageRouter.Route(
			&stream2.FullMarketData_Clear{
				Instrument: msg.Instrument,
			},
		)
	}
	delete(self.fmd, msg.Instrument)
	self.publishListOfInstruments = true
}

func (self *data) handlePublishFullMarketData(msg *PublishFullMarketData) {
	key := outstandingRequests{
		adder: msg.PubSubBag,
		s:     msg.PublishInstrument,
	}
	self.outstandingRequestsMap[key] = msg
}

func (self *data) handleEmptyQueue(msg *messages.EmptyQueue) {
	// Publish List
	if self.publishListOfInstruments {
		ss, _ := self.buildInstrumentList()
		InstrumentStatusArray := make([]*stream2.InstrumentStatus, len(ss), len(ss))
		for i, s := range ss {
			InstrumentStatusArray[i] = &stream2.InstrumentStatus{
				Instrument: s.Instrument,
				Status:     s.Status,
			}
		}
		if self.pubSub.PubWithContext(
			&stream2.FullMarketData_InstrumentList_Response{
				Instruments: InstrumentStatusArray,
			},
			self.fullMarketDataHelper.InstrumentListChannelName(),
		) {
			self.publishListOfInstruments = false
		} else {
			msg.ErrorHappen = true
			return
		}
	}

	for bag, request := range self.outstandingRequestsMap {
		switch req := request.(type) {
		case *PublishFullMarketData:
			if v, ok := self.fmd[req.PublishInstrument]; ok {
				top5, b := self.calculate(true, v)
				if b {
					bag.adder.Add(top5)
				}
			}
		case *stream2.FullMarketData_Instrument_RegisterWrapper:
			self.doFullMarketData_Instrument_RegisterWrapper(req)
		}

	}
	self.outstandingRequestsMap = make(map[outstandingRequests]interface{})

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
			self.fullMarketDataHelper.InstrumentChannelNameForTop5(key),
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

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_AddOrderInstructionWrapper(msg *stream2.FullMarketData_AddOrderInstructionWrapper) {
	self.MessageRouter.Route(msg.Data)
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_AddOrderInstruction(msg *stream2.FullMarketData_AddOrderInstruction) {
	if _, ok := self.fmdCount[msg.Instrument]; ok || !self.proxy {
		if orderBook, ok := self.findFullMarketDataBook(msg.FeedName, msg.Instrument); ok {
			_ = orderBook.Send(msg)
		}
		self.pubSub.Pub(msg, self.fullMarketDataHelper.InstrumentChannelName(msg.Instrument))
		self.fmdDirty[msg.Instrument] = true
	}
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_ClearWrapper(msg *stream2.FullMarketData_ClearWrapper) {
	self.MessageRouter.Route(msg.Data)
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_Clear(msg *stream2.FullMarketData_Clear) {
	if book, ok := self.findFullMarketDataBook(msg.FeedName, msg.Instrument); ok {
		_ = book.Send(msg)
	}
	self.pubSub.Pub(msg, self.fullMarketDataHelper.InstrumentChannelName(msg.Instrument))
	self.fmdDirty[msg.Instrument] = true
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_ReduceVolumeInstructionWrapper(msg *stream2.FullMarketData_ReduceVolumeInstructionWrapper) {
	self.MessageRouter.Route(msg.Data)
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_ReduceVolumeInstruction(msg *stream2.FullMarketData_ReduceVolumeInstruction) {
	if _, ok := self.fmdCount[msg.Instrument]; ok || !self.proxy {
		if book, ok := self.findFullMarketDataBook(msg.FeedName, msg.Instrument); ok {
			_ = book.Send(msg)
		}
		self.pubSub.Pub(msg, self.fullMarketDataHelper.InstrumentChannelName(msg.Instrument))
		self.fmdDirty[msg.Instrument] = true
	}
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_DeleteOrderInstructionWrapper(msg *stream2.FullMarketData_DeleteOrderInstructionWrapper) {
	self.MessageRouter.Route(msg.Data)
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_Instrument_InstrumentStatusWrapper(msg *stream2.FullMarketData_Instrument_InstrumentStatusWrapper) {
	self.MessageRouter.Route(msg.Data)
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_Instrument_InstrumentStatus(msg *stream2.FullMarketData_Instrument_InstrumentStatus) {
	if book, ok := self.findFullMarketDataBook(msg.FeedName, msg.Instrument); ok {
		_ = book.Send(msg)
	}
	self.pubSub.Pub(msg, self.fullMarketDataHelper.InstrumentChannelName(msg.Instrument))
	self.publishListOfInstruments = true
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_DeleteOrderInstruction(msg *stream2.FullMarketData_DeleteOrderInstruction) {
	if _, ok := self.fmdCount[msg.Instrument]; ok || !self.proxy {
		if book, ok := self.findFullMarketDataBook(msg.FeedName, msg.Instrument); ok {
			_ = book.Send(msg)
		}
		self.pubSub.Pub(msg, self.fullMarketDataHelper.InstrumentChannelName(msg.Instrument))
		self.fmdDirty[msg.Instrument] = true
	}
}

func (self *data) calculate(force bool, fullMarketOrderBook fullMarketData.IFullMarketOrderBook) (*stream.PublishTop5, bool) {
	thereWasAChange := force || fullMarketOrderBook.OrderCount() == 0
	var bids []*stream.Point
	if highBidNode := fullMarketOrderBook.BidOrderSide().Right(); highBidNode != nil {
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
	if lowAskNode := fullMarketOrderBook.AskOrderSide().Left(); lowAskNode != nil {
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
			Instrument: fullMarketOrderBook.InstrumentName(),
			Spread:     spread,
			Bid:        bids,
			Ask:        asks,
		}, true
	}
	return nil, false
}

func (self *data) buildInstrumentList() ([]InstrumentStatus, error) {
	var ss []string
	for key := range self.fmd {
		ss = append(ss, key)
	}
	sort.Strings(ss)
	result := make([]InstrumentStatus, len(ss), len(ss))
	for i, s := range ss {
		if status, ok := self.fmd[s]; ok {
			result[i] = InstrumentStatus{
				Instrument: s,
				Status:     status.Status(),
			}
			continue
		}
		result[i] = InstrumentStatus{
			Instrument: s,
			Status:     "",
		}
	}

	return result, nil
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) doFullMarketData_Instrument_RegisterWrapper(request *stream2.FullMarketData_Instrument_RegisterWrapper) {
	if fmdBook, ok := self.fmd[request.Data.Instrument]; ok {
		localSideEncode := func(
			orderSide *avltree.Tree,
			startNodeFn func(orderSide *avltree.Tree) *avltree.Node,
			direction func(node *avltree.Node) *avltree.Node,
		) {
			startNode := startNodeFn(orderSide)
			if firstNode := startNode; firstNode != nil {
				messagesToNext := make([]proto.Message, 0, orderSide.Size())
				for node := firstNode; node != nil; node = direction(node) {
					if pp, ok := node.Value.(*fullMarketData.PricePoint); ok {
						iterator := pp.List.Iterator()
						for iterator.Next() {
							if order, ok := iterator.Value().(*fullMarketData.FullMarketOrder); ok {

								messagesToNext = append(
									messagesToNext,
									&stream2.FullMarketData_AddOrderInstruction{
										Instrument: request.Data.Instrument,
										Order: &stream2.FullMarketData_AddOrder{
											Side:   order.Side,
											Id:     order.Id,
											Price:  order.Price,
											Volume: order.Volume,
										},
									},
								)
							}
						}
					}
				}
				request.ToNext(messagesToNext)
			}
		}

		localSideEncode(
			fmdBook.BidOrderSide(),
			func(orderSide *avltree.Tree) *avltree.Node {
				return orderSide.Right()
			},
			func(node *avltree.Node) *avltree.Node {
				node = node.Prev()
				return node
			},
		)

		localSideEncode(
			fmdBook.AskOrderSide(),
			func(orderSide *avltree.Tree) *avltree.Node {
				return orderSide.Left()
			},
			func(node *avltree.Node) *avltree.Node {
				node = node.Next()
				return node
			},
		)
	}
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_Instrument_RegisterWrapper(msg *stream2.FullMarketData_Instrument_RegisterWrapper) {
	self.SubscribeFullMarketData(msg.Data.Instrument)
	msg.ToNext(
		&stream2.FullMarketData_Clear{
			Instrument: msg.Data.Instrument,
		},
	)
	key := outstandingRequests{
		adder: msg.Adder(),
		s:     msg.Data.Instrument,
	}
	self.outstandingRequestsMap[key] = msg
}

func (self *data) handleRequestAllInstruments(request *RequestAllInstruments) {
	self.doInstrumentListRequest(request.Next)
}

func (self *data) handleCallbackMessage(request *CallbackMessage) {
	if v, ok := self.fmd[request.InstrumentName]; ok {
		if request.CallBack != nil {
			msg := request.CallBack(request.Data, v)
			if msg != nil {
				self.MessageRouter.MultiRoute(msg...)
			}
		}
	}
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_Instrument_UnregisterWrapper(request *stream2.FullMarketData_Instrument_UnregisterWrapper) {
	self.UnsubscribeFullMarketData(request.Data.Instrument)
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_InstrumentList_RequestWrapper(request *stream2.FullMarketData_InstrumentList_RequestWrapper) {
	self.doInstrumentListRequest(request.ToNext)
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) doInstrumentListRequest(cb rxgo.NextFunc) {
	if cb == nil {
		return
	}
	ss, err := self.buildInstrumentList()
	if err != nil {
		return
	}
	array := make([]*stream2.InstrumentStatus, len(ss), len(ss))
	for i, s := range ss {
		array[i] = &stream2.InstrumentStatus{
			Instrument: s.Instrument,
			Status:     s.Status,
		}
	}
	cb(
		&stream2.FullMarketData_InstrumentList_Response{
			Instruments: array,
		},
	)
}

//goland:noinspection GoSnakeCaseUsage
func (self *data) handleFullMarketData_InstrumentList_ResponseWrapper(incomingMessage *stream2.FullMarketData_InstrumentList_ResponseWrapper) {
	m := make(map[string]string)
	for _, instrument := range incomingMessage.Data.Instruments {
		m[instrument.Instrument] = instrument.Status
	}

	feedFmdMap := make(map[string]fullMarketData.IFullMarketOrderBook)
	otherFmdMap := make(map[string]fullMarketData.IFullMarketOrderBook)
	for key, orderBook := range self.fmd {
		if orderBook.FeedName() == incomingMessage.Data.FeedName {
			if v, ok := m[key]; ok {
				orderBook.SetStatus(v)
				feedFmdMap[key] = orderBook
				delete(m, key)
			}
		} else {
			otherFmdMap[key] = orderBook
		}
	}
	for k, v := range m {
		newInstrument := fullMarketData.NewFullMarketOrderBook(incomingMessage.Data.FeedName, k)
		newInstrument.SetStatus(v)
		feedFmdMap[k] = newInstrument
	}

	self.fmd = make(map[string]fullMarketData.IFullMarketOrderBook)
	for key, value := range otherFmdMap {
		self.fmd[key] = value
	}
	for key, value := range feedFmdMap {
		self.fmd[key] = value
	}

	self.publishListOfInstruments = true
}

func newData(pubSub *pubsub.PubSub, fullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper, proxy bool) (IFmdManagerData, error) {
	result := &data{
		proxy:                  proxy,
		outstandingRequestsMap: make(map[outstandingRequests]interface{}),
		queuedMessages:         make(map[string]*stream.PublishTop5),
		MessageRouter:          messageRouter.NewMessageRouter(),
		pubSub:                 pubSub,
		fmd:                    make(map[string]fullMarketData.IFullMarketOrderBook),
		fmdDirty:               make(map[string]bool),
		fmdCount:               make(map[string]int),
		fullMarketDataHelper:   fullMarketDataHelper,
	}
	_ = result.MessageRouter.Add(result.handleEmptyQueue)
	//
	_ = result.MessageRouter.Add(result.handleFullMarketData_AddOrderInstructionWrapper)
	_ = result.MessageRouter.Add(result.handleFullMarketData_ClearWrapper)
	_ = result.MessageRouter.Add(result.handleFullMarketData_ReduceVolumeInstructionWrapper)
	_ = result.MessageRouter.Add(result.handleFullMarketData_DeleteOrderInstructionWrapper)
	_ = result.MessageRouter.Add(result.handleFullMarketData_Instrument_InstrumentStatusWrapper)
	//
	_ = result.MessageRouter.Add(result.handleFullMarketData_AddOrderInstruction)
	_ = result.MessageRouter.Add(result.handleFullMarketData_Clear)
	_ = result.MessageRouter.Add(result.handleFullMarketData_ReduceVolumeInstruction)
	_ = result.MessageRouter.Add(result.handleFullMarketData_DeleteOrderInstruction)
	_ = result.MessageRouter.Add(result.handleFullMarketData_Instrument_InstrumentStatus)
	//
	_ = result.MessageRouter.Add(result.handlePublishFullMarketData)
	_ = result.MessageRouter.Add(result.handleFullMarketDataRemoveInstrumentInstruction)
	//
	_ = result.MessageRouter.Add(result.handleFullMarketData_InstrumentList_RequestWrapper)
	_ = result.MessageRouter.Add(result.handleFullMarketData_InstrumentList_ResponseWrapper)
	//
	_ = result.MessageRouter.Add(result.handleFullMarketData_Instrument_RegisterWrapper)
	_ = result.MessageRouter.Add(result.handleFullMarketData_Instrument_UnregisterWrapper)
	//
	_ = result.MessageRouter.Add(result.handleCallbackMessage)
	_ = result.MessageRouter.Add(result.handleRequestAllInstruments)

	//
	return result, nil
}
