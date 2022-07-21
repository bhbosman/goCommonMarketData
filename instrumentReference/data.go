package instrumentReference

import (
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
)

type data struct {
	MessageRouter *messageRouter.MessageRouter
	m             map[string]*ReferenceData
}

func (self *data) GetReferenceData(instrumentData string) (*ReferenceData, bool) {
	if v, ok := self.m[instrumentData]; ok {
		return v, true
	}
	return nil, false
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

func (self *data) handleEmptyQueue(msg *messages.EmptyQueue) {
}

func (self *data) handleReferenceData(msg *ReferenceData) {
	self.m[msg.Name] = msg
}

func newData() (IInstrumentReferenceData, error) {
	result := &data{
		MessageRouter: messageRouter.NewMessageRouter(),
		m:             make(map[string]*ReferenceData),
	}
	result.MessageRouter.Add(result.handleEmptyQueue)
	result.MessageRouter.Add(result.handleReferenceData)

	return result, nil
}
