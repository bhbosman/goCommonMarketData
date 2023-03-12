package fullMarketDataManagerService

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketData/stream"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/gocommon/pubSub"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataRegistration(t *testing.T) {
	helper, _ := fullMarketDataHelper.NewFullMarketDataHelper()
	factory := func(mockController *gomock.Controller, sub pubSub.IPubSub) *data {

		sut, err := newData(sub, helper, true)
		assert.NoError(t, err)
		assert.NotNil(t, sut)
		return sut.(*data)
	}
	t.Run(
		"Call Register and Deregister for two clients on one instrument",
		func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			instrument := "Instrument01"
			registeredSource := helper.RegisteredSource(instrument)

			sub := pubSub.NewMockIPubSub(mockController)
			sub.EXPECT().Pub(
				&stream.FullMarketData_Instrument_Register{
					Instrument:   instrument,
					RegisterName: "",
				},
				registeredSource).Times(1)
			sub.EXPECT().Pub(
				&stream.FullMarketData_Instrument_Unregister{
					Instrument:   instrument,
					RegisterName: "",
				},
				registeredSource).Times(1)

			sut := factory(mockController, sub)
			sut.SubscribeFullMarketData("Client01", instrument)
			assert.Len(t, sut.fmdCount, 1)
			assert.Equal(t, 1, sut.fmdCount[instrument].Count())
			sut.SubscribeFullMarketData("Client02", instrument)
			assert.Equal(t, 2, sut.fmdCount[instrument].Count())
			sut.UnsubscribeFullMarketData("Client01", instrument)
			assert.Len(t, sut.fmdCount, 1)
			assert.Equal(t, 1, sut.fmdCount[instrument].Count())
			sut.UnsubscribeFullMarketData("Client02", instrument)
			assert.Len(t, sut.fmdCount, 0)
		},
	)
	t.Run(
		"One One instrument, register two clients, and deregister one client twice",
		func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			instrument := "Instrument01"
			registeredSource := helper.RegisteredSource(instrument)

			sub := pubSub.NewMockIPubSub(mockController)
			sub.EXPECT().Pub(
				&stream.FullMarketData_Instrument_Register{
					Instrument:   instrument,
					RegisterName: "",
				},
				registeredSource).Times(1)
			sub.EXPECT().Pub(
				&stream.FullMarketData_Instrument_Unregister{
					Instrument:   instrument,
					RegisterName: "",
				},
				registeredSource).Times(0)

			sut := factory(mockController, sub)
			sut.SubscribeFullMarketData("Client01", instrument)
			assert.Len(t, sut.fmdCount, 1)
			assert.Equal(t, 1, sut.fmdCount[instrument].Count())
			sut.SubscribeFullMarketData("Client02", instrument)
			assert.Equal(t, 2, sut.fmdCount[instrument].Count())
			sut.UnsubscribeFullMarketData("Client01", instrument)
			assert.Len(t, sut.fmdCount, 1)
			assert.Equal(t, 1, sut.fmdCount[instrument].Count())
			sut.UnsubscribeFullMarketData("Client01", instrument)
			assert.Len(t, sut.fmdCount, 1)
			assert.Equal(t, 1, sut.fmdCount[instrument].Count())
		},
	)

}
