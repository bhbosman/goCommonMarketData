// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bhbosman/goCommonMarketData/fullMarketDataManagerViewer (interfaces: IFullMarketDataView)

// Package fullMarketDataManagerViewer is a generated GoMock package.
package fullMarketDataManagerViewer

import (
	fmt "fmt"

	stream "github.com/bhbosman/goMessages/marketData/stream"
	errors "github.com/bhbosman/gocommon/errors"
	"golang.org/x/net/context"
)

// Interface A Comment
// Interface github.com/bhbosman/goCommonMarketData/fullMarketDataManagerViewer
// Interface IFullMarketDataView
// Interface IFullMarketDataView, Method: Send
type IFullMarketDataViewSendIn struct {
	arg0 interface{}
}

type IFullMarketDataViewSendOut struct {
	Args0 error
}
type IFullMarketDataViewSendError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *IFullMarketDataViewSendError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type IFullMarketDataViewSend struct {
	inData         IFullMarketDataViewSendIn
	outDataChannel chan IFullMarketDataViewSendOut
}

func NewIFullMarketDataViewSend(waitToComplete bool, arg0 interface{}) *IFullMarketDataViewSend {
	var outDataChannel chan IFullMarketDataViewSendOut
	if waitToComplete {
		outDataChannel = make(chan IFullMarketDataViewSendOut)
	} else {
		outDataChannel = nil
	}
	return &IFullMarketDataViewSend{
		inData: IFullMarketDataViewSendIn{
			arg0: arg0,
		},
		outDataChannel: outDataChannel,
	}
}

func (self *IFullMarketDataViewSend) Wait(onError func(interfaceName string, methodName string, err error) error) (IFullMarketDataViewSendOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &IFullMarketDataViewSendError{
			InterfaceName: "IFullMarketDataView",
			MethodName:    "Send",
			Reason:        "Channel for IFullMarketDataView::Send returned false",
		}
		if onError != nil {
			err := onError("IFullMarketDataView", "Send", generatedError)
			return IFullMarketDataViewSendOut{}, err
		} else {
			return IFullMarketDataViewSendOut{}, generatedError
		}
	}
	return data, nil
}

func (self *IFullMarketDataViewSend) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallIFullMarketDataViewSend(context context.Context, channel chan<- interface{}, waitToComplete bool, arg0 interface{}) (IFullMarketDataViewSendOut, error) {
	if context != nil && context.Err() != nil {
		return IFullMarketDataViewSendOut{}, context.Err()
	}
	data := NewIFullMarketDataViewSend(waitToComplete, arg0)
	if waitToComplete {
		defer func(data *IFullMarketDataViewSend) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return IFullMarketDataViewSendOut{}, context.Err()
	}
	channel <- data
	var err error
	var v IFullMarketDataViewSendOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return IFullMarketDataViewSendOut{}, err
	}
	return v, nil
}

// Interface IFullMarketDataView, Method: SetMarketDataInstanceChange
type IFullMarketDataViewSetMarketDataInstanceChangeIn struct {
	arg0 func(*stream.PublishTop5)
}

type IFullMarketDataViewSetMarketDataInstanceChangeOut struct {
}
type IFullMarketDataViewSetMarketDataInstanceChangeError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *IFullMarketDataViewSetMarketDataInstanceChangeError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type IFullMarketDataViewSetMarketDataInstanceChange struct {
	inData         IFullMarketDataViewSetMarketDataInstanceChangeIn
	outDataChannel chan IFullMarketDataViewSetMarketDataInstanceChangeOut
}

func NewIFullMarketDataViewSetMarketDataInstanceChange(waitToComplete bool, arg0 func(*stream.PublishTop5)) *IFullMarketDataViewSetMarketDataInstanceChange {
	var outDataChannel chan IFullMarketDataViewSetMarketDataInstanceChangeOut
	if waitToComplete {
		outDataChannel = make(chan IFullMarketDataViewSetMarketDataInstanceChangeOut)
	} else {
		outDataChannel = nil
	}
	return &IFullMarketDataViewSetMarketDataInstanceChange{
		inData: IFullMarketDataViewSetMarketDataInstanceChangeIn{
			arg0: arg0,
		},
		outDataChannel: outDataChannel,
	}
}

func (self *IFullMarketDataViewSetMarketDataInstanceChange) Wait(onError func(interfaceName string, methodName string, err error) error) (IFullMarketDataViewSetMarketDataInstanceChangeOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &IFullMarketDataViewSetMarketDataInstanceChangeError{
			InterfaceName: "IFullMarketDataView",
			MethodName:    "SetMarketDataInstanceChange",
			Reason:        "Channel for IFullMarketDataView::SetMarketDataInstanceChange returned false",
		}
		if onError != nil {
			err := onError("IFullMarketDataView", "SetMarketDataInstanceChange", generatedError)
			return IFullMarketDataViewSetMarketDataInstanceChangeOut{}, err
		} else {
			return IFullMarketDataViewSetMarketDataInstanceChangeOut{}, generatedError
		}
	}
	return data, nil
}

func (self *IFullMarketDataViewSetMarketDataInstanceChange) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallIFullMarketDataViewSetMarketDataInstanceChange(context context.Context, channel chan<- interface{}, waitToComplete bool, arg0 func(*stream.PublishTop5)) (IFullMarketDataViewSetMarketDataInstanceChangeOut, error) {
	if context != nil && context.Err() != nil {
		return IFullMarketDataViewSetMarketDataInstanceChangeOut{}, context.Err()
	}
	data := NewIFullMarketDataViewSetMarketDataInstanceChange(waitToComplete, arg0)
	if waitToComplete {
		defer func(data *IFullMarketDataViewSetMarketDataInstanceChange) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return IFullMarketDataViewSetMarketDataInstanceChangeOut{}, context.Err()
	}
	channel <- data
	var err error
	var v IFullMarketDataViewSetMarketDataInstanceChangeOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return IFullMarketDataViewSetMarketDataInstanceChangeOut{}, err
	}
	return v, nil
}

// Interface IFullMarketDataView, Method: SetMarketDataListChange
type IFullMarketDataViewSetMarketDataListChangeIn struct {
	arg0 func([]string)
}

type IFullMarketDataViewSetMarketDataListChangeOut struct {
}
type IFullMarketDataViewSetMarketDataListChangeError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *IFullMarketDataViewSetMarketDataListChangeError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type IFullMarketDataViewSetMarketDataListChange struct {
	inData         IFullMarketDataViewSetMarketDataListChangeIn
	outDataChannel chan IFullMarketDataViewSetMarketDataListChangeOut
}

func NewIFullMarketDataViewSetMarketDataListChange(waitToComplete bool, arg0 func([]string)) *IFullMarketDataViewSetMarketDataListChange {
	var outDataChannel chan IFullMarketDataViewSetMarketDataListChangeOut
	if waitToComplete {
		outDataChannel = make(chan IFullMarketDataViewSetMarketDataListChangeOut)
	} else {
		outDataChannel = nil
	}
	return &IFullMarketDataViewSetMarketDataListChange{
		inData: IFullMarketDataViewSetMarketDataListChangeIn{
			arg0: arg0,
		},
		outDataChannel: outDataChannel,
	}
}

func (self *IFullMarketDataViewSetMarketDataListChange) Wait(onError func(interfaceName string, methodName string, err error) error) (IFullMarketDataViewSetMarketDataListChangeOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &IFullMarketDataViewSetMarketDataListChangeError{
			InterfaceName: "IFullMarketDataView",
			MethodName:    "SetMarketDataListChange",
			Reason:        "Channel for IFullMarketDataView::SetMarketDataListChange returned false",
		}
		if onError != nil {
			err := onError("IFullMarketDataView", "SetMarketDataListChange", generatedError)
			return IFullMarketDataViewSetMarketDataListChangeOut{}, err
		} else {
			return IFullMarketDataViewSetMarketDataListChangeOut{}, generatedError
		}
	}
	return data, nil
}

func (self *IFullMarketDataViewSetMarketDataListChange) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallIFullMarketDataViewSetMarketDataListChange(context context.Context, channel chan<- interface{}, waitToComplete bool, arg0 func([]string)) (IFullMarketDataViewSetMarketDataListChangeOut, error) {
	if context != nil && context.Err() != nil {
		return IFullMarketDataViewSetMarketDataListChangeOut{}, context.Err()
	}
	data := NewIFullMarketDataViewSetMarketDataListChange(waitToComplete, arg0)
	if waitToComplete {
		defer func(data *IFullMarketDataViewSetMarketDataListChange) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return IFullMarketDataViewSetMarketDataListChangeOut{}, context.Err()
	}
	channel <- data
	var err error
	var v IFullMarketDataViewSetMarketDataListChangeOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return IFullMarketDataViewSetMarketDataListChangeOut{}, err
	}
	return v, nil
}

// Interface IFullMarketDataView, Method: SubscribeFullMarketData
type IFullMarketDataViewSubscribeFullMarketDataIn struct {
	arg0 string
}

type IFullMarketDataViewSubscribeFullMarketDataOut struct {
}
type IFullMarketDataViewSubscribeFullMarketDataError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *IFullMarketDataViewSubscribeFullMarketDataError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type IFullMarketDataViewSubscribeFullMarketData struct {
	inData         IFullMarketDataViewSubscribeFullMarketDataIn
	outDataChannel chan IFullMarketDataViewSubscribeFullMarketDataOut
}

func NewIFullMarketDataViewSubscribeFullMarketData(waitToComplete bool, arg0 string) *IFullMarketDataViewSubscribeFullMarketData {
	var outDataChannel chan IFullMarketDataViewSubscribeFullMarketDataOut
	if waitToComplete {
		outDataChannel = make(chan IFullMarketDataViewSubscribeFullMarketDataOut)
	} else {
		outDataChannel = nil
	}
	return &IFullMarketDataViewSubscribeFullMarketData{
		inData: IFullMarketDataViewSubscribeFullMarketDataIn{
			arg0: arg0,
		},
		outDataChannel: outDataChannel,
	}
}

func (self *IFullMarketDataViewSubscribeFullMarketData) Wait(onError func(interfaceName string, methodName string, err error) error) (IFullMarketDataViewSubscribeFullMarketDataOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &IFullMarketDataViewSubscribeFullMarketDataError{
			InterfaceName: "IFullMarketDataView",
			MethodName:    "SubscribeFullMarketData",
			Reason:        "Channel for IFullMarketDataView::SubscribeFullMarketData returned false",
		}
		if onError != nil {
			err := onError("IFullMarketDataView", "SubscribeFullMarketData", generatedError)
			return IFullMarketDataViewSubscribeFullMarketDataOut{}, err
		} else {
			return IFullMarketDataViewSubscribeFullMarketDataOut{}, generatedError
		}
	}
	return data, nil
}

func (self *IFullMarketDataViewSubscribeFullMarketData) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallIFullMarketDataViewSubscribeFullMarketData(context context.Context, channel chan<- interface{}, waitToComplete bool, arg0 string) (IFullMarketDataViewSubscribeFullMarketDataOut, error) {
	if context != nil && context.Err() != nil {
		return IFullMarketDataViewSubscribeFullMarketDataOut{}, context.Err()
	}
	data := NewIFullMarketDataViewSubscribeFullMarketData(waitToComplete, arg0)
	if waitToComplete {
		defer func(data *IFullMarketDataViewSubscribeFullMarketData) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return IFullMarketDataViewSubscribeFullMarketDataOut{}, context.Err()
	}
	channel <- data
	var err error
	var v IFullMarketDataViewSubscribeFullMarketDataOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return IFullMarketDataViewSubscribeFullMarketDataOut{}, err
	}
	return v, nil
}

// Interface IFullMarketDataView, Method: UnsubscribeFullMarketData
type IFullMarketDataViewUnsubscribeFullMarketDataIn struct {
	arg0 string
}

type IFullMarketDataViewUnsubscribeFullMarketDataOut struct {
}
type IFullMarketDataViewUnsubscribeFullMarketDataError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *IFullMarketDataViewUnsubscribeFullMarketDataError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type IFullMarketDataViewUnsubscribeFullMarketData struct {
	inData         IFullMarketDataViewUnsubscribeFullMarketDataIn
	outDataChannel chan IFullMarketDataViewUnsubscribeFullMarketDataOut
}

func NewIFullMarketDataViewUnsubscribeFullMarketData(waitToComplete bool, arg0 string) *IFullMarketDataViewUnsubscribeFullMarketData {
	var outDataChannel chan IFullMarketDataViewUnsubscribeFullMarketDataOut
	if waitToComplete {
		outDataChannel = make(chan IFullMarketDataViewUnsubscribeFullMarketDataOut)
	} else {
		outDataChannel = nil
	}
	return &IFullMarketDataViewUnsubscribeFullMarketData{
		inData: IFullMarketDataViewUnsubscribeFullMarketDataIn{
			arg0: arg0,
		},
		outDataChannel: outDataChannel,
	}
}

func (self *IFullMarketDataViewUnsubscribeFullMarketData) Wait(onError func(interfaceName string, methodName string, err error) error) (IFullMarketDataViewUnsubscribeFullMarketDataOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &IFullMarketDataViewUnsubscribeFullMarketDataError{
			InterfaceName: "IFullMarketDataView",
			MethodName:    "UnsubscribeFullMarketData",
			Reason:        "Channel for IFullMarketDataView::UnsubscribeFullMarketData returned false",
		}
		if onError != nil {
			err := onError("IFullMarketDataView", "UnsubscribeFullMarketData", generatedError)
			return IFullMarketDataViewUnsubscribeFullMarketDataOut{}, err
		} else {
			return IFullMarketDataViewUnsubscribeFullMarketDataOut{}, generatedError
		}
	}
	return data, nil
}

func (self *IFullMarketDataViewUnsubscribeFullMarketData) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallIFullMarketDataViewUnsubscribeFullMarketData(context context.Context, channel chan<- interface{}, waitToComplete bool, arg0 string) (IFullMarketDataViewUnsubscribeFullMarketDataOut, error) {
	if context != nil && context.Err() != nil {
		return IFullMarketDataViewUnsubscribeFullMarketDataOut{}, context.Err()
	}
	data := NewIFullMarketDataViewUnsubscribeFullMarketData(waitToComplete, arg0)
	if waitToComplete {
		defer func(data *IFullMarketDataViewUnsubscribeFullMarketData) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return IFullMarketDataViewUnsubscribeFullMarketDataOut{}, context.Err()
	}
	channel <- data
	var err error
	var v IFullMarketDataViewUnsubscribeFullMarketDataOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return IFullMarketDataViewUnsubscribeFullMarketDataOut{}, err
	}
	return v, nil
}

func ChannelEventsForIFullMarketDataView(next IFullMarketDataView, event interface{}) (bool, error) {
	switch v := event.(type) {
	case *IFullMarketDataViewSend:
		data := IFullMarketDataViewSendOut{}
		data.Args0 = next.Send(v.inData.arg0)
		if v.outDataChannel != nil {
			v.outDataChannel <- data
		}
		return true, nil
	case *IFullMarketDataViewSubscribeFullMarketData:
		data := IFullMarketDataViewSubscribeFullMarketDataOut{}
		next.SubscribeFullMarketData(v.inData.arg0)
		if v.outDataChannel != nil {
			v.outDataChannel <- data
		}
		return true, nil
	case *IFullMarketDataViewUnsubscribeFullMarketData:
		data := IFullMarketDataViewUnsubscribeFullMarketDataOut{}
		next.UnsubscribeFullMarketData(v.inData.arg0)
		if v.outDataChannel != nil {
			v.outDataChannel <- data
		}
		return true, nil
	default:
		return false, nil
	}
}
