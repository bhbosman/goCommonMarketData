package instrumentReference

import (
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
)

type data struct {
	MessageRouter   *messageRouter.MessageRouter
	lunoProviders   map[string]*LunoReferenceData
	referenceData   map[string]IReferenceData
	krakenProviders map[string]*KrakenReferenceData
}

func (self *data) GetKrakenProviders() ([]KrakenReferenceData, error) {
	result := make([]KrakenReferenceData, len(self.krakenProviders), len(self.krakenProviders))
	i := 0
	for _, value := range self.krakenProviders {
		result[i] = *value
		i++
	}
	return result, nil
}

func (self *data) GetLunoProviders() ([]LunoReferenceData, error) {
	result := make([]LunoReferenceData, len(self.lunoProviders), len(self.lunoProviders))
	i := 0
	for _, value := range self.lunoProviders {
		result[i] = *value
		i++
	}
	return result, nil
}

func (self *data) GetReferenceData(instrumentData string) (IReferenceData, bool) {
	if v, ok := self.referenceData[instrumentData]; ok {
		return v, true
	}
	return nil, false
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

func (self *data) handleEmptyQueue(msg *messages.EmptyQueue) {
}

func (self *data) handleMarketDataFeedReference(msg *MarketDataFeedReference) {
	for _, feed := range msg.LunoFeeds {
		self.referenceData[feed.SystemName] = feed
		self.lunoProviders[feed.SystemName] = feed
	}
	for _, feed := range msg.KrakenFeeds {
		self.krakenProviders[feed.ConnectionName] = feed
		for _, krakenFeed := range feed.Feeds {
			self.referenceData[krakenFeed.SystemName] = krakenFeed
		}
	}
}

func newData() (IInstrumentReferenceData, error) {
	result := &data{
		MessageRouter:   messageRouter.NewMessageRouter(),
		lunoProviders:   make(map[string]*LunoReferenceData),
		referenceData:   make(map[string]IReferenceData),
		krakenProviders: make(map[string]*KrakenReferenceData),
	}
	result.MessageRouter.Add(result.handleEmptyQueue)
	result.MessageRouter.Add(result.handleMarketDataFeedReference)

	return result, nil
}
