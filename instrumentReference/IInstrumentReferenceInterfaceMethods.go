// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bhbosman/goCommonMarketData/instrumentReference (interfaces: IInstrumentReference)

// Package instrumentReference is a generated GoMock package.
package instrumentReference

import (
	"context"
	fmt "fmt"

	errors "github.com/bhbosman/gocommon/errors"
)

// Interface A Comment
// Interface github.com/bhbosman/goCommonMarketData/instrumentReference
// Interface IInstrumentReference
// Interface IInstrumentReference, Method: GetKrakenProviders
type IInstrumentReferenceGetKrakenProvidersIn struct {
}

type IInstrumentReferenceGetKrakenProvidersOut struct {
	Args0 []KrakenReferenceData
	Args1 error
}
type IInstrumentReferenceGetKrakenProvidersError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *IInstrumentReferenceGetKrakenProvidersError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type IInstrumentReferenceGetKrakenProviders struct {
	inData         IInstrumentReferenceGetKrakenProvidersIn
	outDataChannel chan IInstrumentReferenceGetKrakenProvidersOut
}

func NewIInstrumentReferenceGetKrakenProviders(waitToComplete bool) *IInstrumentReferenceGetKrakenProviders {
	var outDataChannel chan IInstrumentReferenceGetKrakenProvidersOut
	if waitToComplete {
		outDataChannel = make(chan IInstrumentReferenceGetKrakenProvidersOut)
	} else {
		outDataChannel = nil
	}
	return &IInstrumentReferenceGetKrakenProviders{
		inData:         IInstrumentReferenceGetKrakenProvidersIn{},
		outDataChannel: outDataChannel,
	}
}

func (self *IInstrumentReferenceGetKrakenProviders) Wait(onError func(interfaceName string, methodName string, err error) error) (IInstrumentReferenceGetKrakenProvidersOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &IInstrumentReferenceGetKrakenProvidersError{
			InterfaceName: "IInstrumentReference",
			MethodName:    "GetKrakenProviders",
			Reason:        "Channel for IInstrumentReference::GetKrakenProviders returned false",
		}
		if onError != nil {
			err := onError("IInstrumentReference", "GetKrakenProviders", generatedError)
			return IInstrumentReferenceGetKrakenProvidersOut{}, err
		} else {
			return IInstrumentReferenceGetKrakenProvidersOut{}, generatedError
		}
	}
	return data, nil
}

func (self *IInstrumentReferenceGetKrakenProviders) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallIInstrumentReferenceGetKrakenProviders(context context.Context, channel chan<- interface{}, waitToComplete bool) (IInstrumentReferenceGetKrakenProvidersOut, error) {
	if context != nil && context.Err() != nil {
		return IInstrumentReferenceGetKrakenProvidersOut{}, context.Err()
	}
	data := NewIInstrumentReferenceGetKrakenProviders(waitToComplete)
	if waitToComplete {
		defer func(data *IInstrumentReferenceGetKrakenProviders) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return IInstrumentReferenceGetKrakenProvidersOut{}, context.Err()
	}
	channel <- data
	var err error
	var v IInstrumentReferenceGetKrakenProvidersOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return IInstrumentReferenceGetKrakenProvidersOut{}, err
	}
	return v, nil
}

// Interface IInstrumentReference, Method: GetLunoProviders
type IInstrumentReferenceGetLunoProvidersIn struct {
}

type IInstrumentReferenceGetLunoProvidersOut struct {
	Args0 []LunoReferenceData
	Args1 error
}
type IInstrumentReferenceGetLunoProvidersError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *IInstrumentReferenceGetLunoProvidersError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type IInstrumentReferenceGetLunoProviders struct {
	inData         IInstrumentReferenceGetLunoProvidersIn
	outDataChannel chan IInstrumentReferenceGetLunoProvidersOut
}

func NewIInstrumentReferenceGetLunoProviders(waitToComplete bool) *IInstrumentReferenceGetLunoProviders {
	var outDataChannel chan IInstrumentReferenceGetLunoProvidersOut
	if waitToComplete {
		outDataChannel = make(chan IInstrumentReferenceGetLunoProvidersOut)
	} else {
		outDataChannel = nil
	}
	return &IInstrumentReferenceGetLunoProviders{
		inData:         IInstrumentReferenceGetLunoProvidersIn{},
		outDataChannel: outDataChannel,
	}
}

func (self *IInstrumentReferenceGetLunoProviders) Wait(onError func(interfaceName string, methodName string, err error) error) (IInstrumentReferenceGetLunoProvidersOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &IInstrumentReferenceGetLunoProvidersError{
			InterfaceName: "IInstrumentReference",
			MethodName:    "GetLunoProviders",
			Reason:        "Channel for IInstrumentReference::GetLunoProviders returned false",
		}
		if onError != nil {
			err := onError("IInstrumentReference", "GetLunoProviders", generatedError)
			return IInstrumentReferenceGetLunoProvidersOut{}, err
		} else {
			return IInstrumentReferenceGetLunoProvidersOut{}, generatedError
		}
	}
	return data, nil
}

func (self *IInstrumentReferenceGetLunoProviders) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallIInstrumentReferenceGetLunoProviders(context context.Context, channel chan<- interface{}, waitToComplete bool) (IInstrumentReferenceGetLunoProvidersOut, error) {
	if context != nil && context.Err() != nil {
		return IInstrumentReferenceGetLunoProvidersOut{}, context.Err()
	}
	data := NewIInstrumentReferenceGetLunoProviders(waitToComplete)
	if waitToComplete {
		defer func(data *IInstrumentReferenceGetLunoProviders) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return IInstrumentReferenceGetLunoProvidersOut{}, context.Err()
	}
	channel <- data
	var err error
	var v IInstrumentReferenceGetLunoProvidersOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return IInstrumentReferenceGetLunoProvidersOut{}, err
	}
	return v, nil
}

// Interface IInstrumentReference, Method: GetReferenceData
type IInstrumentReferenceGetReferenceDataIn struct {
	arg0 string
}

type IInstrumentReferenceGetReferenceDataOut struct {
	Args0 IReferenceData
	Args1 bool
}
type IInstrumentReferenceGetReferenceDataError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *IInstrumentReferenceGetReferenceDataError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type IInstrumentReferenceGetReferenceData struct {
	inData         IInstrumentReferenceGetReferenceDataIn
	outDataChannel chan IInstrumentReferenceGetReferenceDataOut
}

func NewIInstrumentReferenceGetReferenceData(waitToComplete bool, arg0 string) *IInstrumentReferenceGetReferenceData {
	var outDataChannel chan IInstrumentReferenceGetReferenceDataOut
	if waitToComplete {
		outDataChannel = make(chan IInstrumentReferenceGetReferenceDataOut)
	} else {
		outDataChannel = nil
	}
	return &IInstrumentReferenceGetReferenceData{
		inData: IInstrumentReferenceGetReferenceDataIn{
			arg0: arg0,
		},
		outDataChannel: outDataChannel,
	}
}

func (self *IInstrumentReferenceGetReferenceData) Wait(onError func(interfaceName string, methodName string, err error) error) (IInstrumentReferenceGetReferenceDataOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &IInstrumentReferenceGetReferenceDataError{
			InterfaceName: "IInstrumentReference",
			MethodName:    "GetReferenceData",
			Reason:        "Channel for IInstrumentReference::GetReferenceData returned false",
		}
		if onError != nil {
			err := onError("IInstrumentReference", "GetReferenceData", generatedError)
			return IInstrumentReferenceGetReferenceDataOut{}, err
		} else {
			return IInstrumentReferenceGetReferenceDataOut{}, generatedError
		}
	}
	return data, nil
}

func (self *IInstrumentReferenceGetReferenceData) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallIInstrumentReferenceGetReferenceData(context context.Context, channel chan<- interface{}, waitToComplete bool, arg0 string) (IInstrumentReferenceGetReferenceDataOut, error) {
	if context != nil && context.Err() != nil {
		return IInstrumentReferenceGetReferenceDataOut{}, context.Err()
	}
	data := NewIInstrumentReferenceGetReferenceData(waitToComplete, arg0)
	if waitToComplete {
		defer func(data *IInstrumentReferenceGetReferenceData) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return IInstrumentReferenceGetReferenceDataOut{}, context.Err()
	}
	channel <- data
	var err error
	var v IInstrumentReferenceGetReferenceDataOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return IInstrumentReferenceGetReferenceDataOut{}, err
	}
	return v, nil
}

// Interface IInstrumentReference, Method: Send
type IInstrumentReferenceSendIn struct {
	arg0 interface{}
}

type IInstrumentReferenceSendOut struct {
	Args0 error
}
type IInstrumentReferenceSendError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *IInstrumentReferenceSendError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type IInstrumentReferenceSend struct {
	inData         IInstrumentReferenceSendIn
	outDataChannel chan IInstrumentReferenceSendOut
}

func NewIInstrumentReferenceSend(waitToComplete bool, arg0 interface{}) *IInstrumentReferenceSend {
	var outDataChannel chan IInstrumentReferenceSendOut
	if waitToComplete {
		outDataChannel = make(chan IInstrumentReferenceSendOut)
	} else {
		outDataChannel = nil
	}
	return &IInstrumentReferenceSend{
		inData: IInstrumentReferenceSendIn{
			arg0: arg0,
		},
		outDataChannel: outDataChannel,
	}
}

func (self *IInstrumentReferenceSend) Wait(onError func(interfaceName string, methodName string, err error) error) (IInstrumentReferenceSendOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &IInstrumentReferenceSendError{
			InterfaceName: "IInstrumentReference",
			MethodName:    "Send",
			Reason:        "Channel for IInstrumentReference::Send returned false",
		}
		if onError != nil {
			err := onError("IInstrumentReference", "Send", generatedError)
			return IInstrumentReferenceSendOut{}, err
		} else {
			return IInstrumentReferenceSendOut{}, generatedError
		}
	}
	return data, nil
}

func (self *IInstrumentReferenceSend) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallIInstrumentReferenceSend(context context.Context, channel chan<- interface{}, waitToComplete bool, arg0 interface{}) (IInstrumentReferenceSendOut, error) {
	if context != nil && context.Err() != nil {
		return IInstrumentReferenceSendOut{}, context.Err()
	}
	data := NewIInstrumentReferenceSend(waitToComplete, arg0)
	if waitToComplete {
		defer func(data *IInstrumentReferenceSend) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return IInstrumentReferenceSendOut{}, context.Err()
	}
	channel <- data
	var err error
	var v IInstrumentReferenceSendOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return IInstrumentReferenceSendOut{}, err
	}
	return v, nil
}

// Interface IInstrumentReference, Method: SomeMethod
type IInstrumentReferenceSomeMethodIn struct {
}

type IInstrumentReferenceSomeMethodOut struct {
}
type IInstrumentReferenceSomeMethodError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *IInstrumentReferenceSomeMethodError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type IInstrumentReferenceSomeMethod struct {
	inData         IInstrumentReferenceSomeMethodIn
	outDataChannel chan IInstrumentReferenceSomeMethodOut
}

func NewIInstrumentReferenceSomeMethod(waitToComplete bool) *IInstrumentReferenceSomeMethod {
	var outDataChannel chan IInstrumentReferenceSomeMethodOut
	if waitToComplete {
		outDataChannel = make(chan IInstrumentReferenceSomeMethodOut)
	} else {
		outDataChannel = nil
	}
	return &IInstrumentReferenceSomeMethod{
		inData:         IInstrumentReferenceSomeMethodIn{},
		outDataChannel: outDataChannel,
	}
}

func (self *IInstrumentReferenceSomeMethod) Wait(onError func(interfaceName string, methodName string, err error) error) (IInstrumentReferenceSomeMethodOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &IInstrumentReferenceSomeMethodError{
			InterfaceName: "IInstrumentReference",
			MethodName:    "SomeMethod",
			Reason:        "Channel for IInstrumentReference::SomeMethod returned false",
		}
		if onError != nil {
			err := onError("IInstrumentReference", "SomeMethod", generatedError)
			return IInstrumentReferenceSomeMethodOut{}, err
		} else {
			return IInstrumentReferenceSomeMethodOut{}, generatedError
		}
	}
	return data, nil
}

func (self *IInstrumentReferenceSomeMethod) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallIInstrumentReferenceSomeMethod(context context.Context, channel chan<- interface{}, waitToComplete bool) (IInstrumentReferenceSomeMethodOut, error) {
	if context != nil && context.Err() != nil {
		return IInstrumentReferenceSomeMethodOut{}, context.Err()
	}
	data := NewIInstrumentReferenceSomeMethod(waitToComplete)
	if waitToComplete {
		defer func(data *IInstrumentReferenceSomeMethod) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return IInstrumentReferenceSomeMethodOut{}, context.Err()
	}
	channel <- data
	var err error
	var v IInstrumentReferenceSomeMethodOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return IInstrumentReferenceSomeMethodOut{}, err
	}
	return v, nil
}

func ChannelEventsForIInstrumentReference(next IInstrumentReference, event interface{}) (bool, error) {
	switch v := event.(type) {
	case *IInstrumentReferenceGetKrakenProviders:
		data := IInstrumentReferenceGetKrakenProvidersOut{}
		data.Args0, data.Args1 = next.GetKrakenProviders()
		if v.outDataChannel != nil {
			v.outDataChannel <- data
		}
		return true, nil
	case *IInstrumentReferenceGetLunoProviders:
		data := IInstrumentReferenceGetLunoProvidersOut{}
		data.Args0, data.Args1 = next.GetLunoProviders()
		if v.outDataChannel != nil {
			v.outDataChannel <- data
		}
		return true, nil
	case *IInstrumentReferenceGetReferenceData:
		data := IInstrumentReferenceGetReferenceDataOut{}
		data.Args0, data.Args1 = next.GetReferenceData(v.inData.arg0)
		if v.outDataChannel != nil {
			v.outDataChannel <- data
		}
		return true, nil
	case *IInstrumentReferenceSend:
		data := IInstrumentReferenceSendOut{}
		data.Args0 = next.Send(v.inData.arg0)
		if v.outDataChannel != nil {
			v.outDataChannel <- data
		}
		return true, nil
	case *IInstrumentReferenceSomeMethod:
		data := IInstrumentReferenceSomeMethodOut{}
		next.SomeMethod()
		if v.outDataChannel != nil {
			v.outDataChannel <- data
		}
		return true, nil
	default:
		return false, nil
	}
}
